package monitor

import (
	"context"
	"errors"
	"io"
	"math/big"
	"strconv"
	"testing"
	"time"

	"log/slog"

	"github.com/stretchr/testify/assert"

	validatorrouter "github.com/primev/mev-commit/contracts-abi/clients/ValidatorOptInRouter"
	"github.com/primev/mev-commit/tools/validators-monitor/api"
	"github.com/primev/mev-commit/tools/validators-monitor/config"
	"github.com/primev/mev-commit/tools/validators-monitor/database"
)

type fakeBeacon struct {
	resp *api.ProposerDutiesResponse
	err  error
}

func (f *fakeBeacon) GetProposerDuties(ctx context.Context, epoch uint64) (*api.ProposerDutiesResponse, error) {
	return f.resp, f.err
}
func (f *fakeBeacon) GetBlockBySlot(ctx context.Context, slot uint64) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return strconv.FormatUint(slot*10, 10), nil
}

type fakeBeaconSlots struct{}

func (f *fakeBeaconSlots) GetProposerDuties(ctx context.Context, epoch uint64) (*api.ProposerDutiesResponse, error) {
	return nil, nil
}
func (f *fakeBeaconSlots) GetBlockBySlot(ctx context.Context, slot uint64) (string, error) {
	if slot == 2 {
		return "", errors.New("no block")
	}
	return strconv.FormatUint(slot*10, 10), nil
}

type fakeOptIn struct {
	statuses []validatorrouter.IValidatorOptInRouterOptInStatus
	err      error
}

func (f *fakeOptIn) CheckValidatorsOptedIn(ctx context.Context, pubkeys []string) ([]validatorrouter.IValidatorOptInRouterOptInStatus, error) {
	return f.statuses, f.err
}

type fakeRelay struct {
	results map[string]api.RelayResult
}

func (f *fakeRelay) QueryRelayData(ctx context.Context, blockNumber uint64) map[string]api.RelayResult {
	return f.results
}

type fakeDashboard struct {
	resp        *api.DashboardResponse
	err         error
	commitments []api.CommitmentData
	commitErr   error
}

func (f *fakeDashboard) GetBlockInfo(ctx context.Context, blockNumber uint64) (*api.DashboardResponse, error) {
	return f.resp, f.err
}

func (f *fakeDashboard) GetCommitmentsByBlock(ctx context.Context, blockNumber uint64) ([]api.CommitmentData, error) {
	return f.commitments, f.commitErr
}

type fakeNotifier struct {
	called bool
	err    error
}

func (f *fakeNotifier) NotifyRelayData(
	ctx context.Context,
	pubkey string,
	index uint64,
	blockNumber uint64,
	slot uint64,
	mevReward *big.Int,
	feeRecipient string,
	relaysWithData []string,
	totalRelays []string,
	blockInfo *api.DashboardResponse,
) error {
	f.called = true
	return f.err
}

type fakeDB struct {
	saved       []*database.RelayRecord
	commitments []*database.CommitmentRecord
	err         error
	commitErr   error
}

func (f *fakeDB) SaveRelayData(ctx context.Context, record *database.RelayRecord) error {
	if f.err != nil {
		return f.err
	}
	f.saved = append(f.saved, record)
	return nil
}
func (f *fakeDB) SaveBlockCommitments(ctx context.Context, commitments []*database.CommitmentRecord) error {
	if f.commitErr != nil {
		return f.commitErr
	}
	f.commitments = append(f.commitments, commitments...)
	return nil
}
func (f *fakeDB) InitSchema(ctx context.Context) error { return nil }
func (f *fakeDB) Close() error                         { return nil }

type fakeCalc struct {
	curEpoch uint64
	toFetch  []uint64
}

func (f *fakeCalc) CurrentSlot() uint64                   { return 0 }
func (f *fakeCalc) CurrentEpoch() uint64                  { return f.curEpoch }
func (f *fakeCalc) TimeUntilNextEpoch() time.Duration     { return 0 }
func (f *fakeCalc) EpochStartTime(epoch uint64) time.Time { return time.Now() }
func (f *fakeCalc) TargetEpoch() uint64                   { return f.curEpoch }
func (f *fakeCalc) EpochsToFetch() []uint64               { return f.toFetch }
func (f *fakeCalc) SlotToEpoch(slot uint64) uint64        { return 0 }

func makeTestMonitor() *DutyMonitor {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{
		FetchIntervalSec: 1,
		RelayURLs:        []string{"r1", "r2"},
		DB:               config.DBConfig{Enabled: false},
	}
	return &DutyMonitor{
		logger:          logger,
		config:          cfg,
		calculator:      &fakeCalc{},
		beacon:          &fakeBeacon{},
		relay:           &fakeRelay{},
		dashboard:       &fakeDashboard{},
		notifier:        &fakeNotifier{},
		optChecker:      &fakeOptIn{},
		dutiesCache:     make(map[uint64]cachedDuties),
		processedBlocks: make(map[uint64]time.Time),
		db:              nil,
	}
}

func TestCheckEpochTransition(t *testing.T) {
	m := makeTestMonitor()
	fc := &fakeCalc{curEpoch: 5}
	m.calculator = fc
	m.runningEpoch = 5

	m.checkEpochTransition() // no change
	assert.Equal(t, uint64(5), m.runningEpoch)

	fc.curEpoch = 6
	m.checkEpochTransition()
	assert.Equal(t, uint64(6), m.runningEpoch)
}

func TestCleanupCaches(t *testing.T) {
	m := makeTestMonitor()

	// add stale duties
	for i := uint64(1); i <= 3; i++ {
		m.dutiesCache[i] = cachedDuties{cachedAt: time.Now().Add(-7 * time.Hour)}
	}
	// add fresh duty
	m.dutiesCache[99] = cachedDuties{cachedAt: time.Now()}

	// add many processed blocks
	for i := uint64(0); i < processedBlocksTarget+500; i++ {
		m.processedBlocks[i] = time.Now().Add(-3 * time.Hour)
	}

	m.cleanupCaches()

	assert.False(t, func() bool { _, ok := m.dutiesCache[1]; return ok }())
	assert.True(t, func() bool { _, ok := m.dutiesCache[99]; return ok }())
	assert.LessOrEqual(t, len(m.processedBlocks), processedBlocksTarget)
}

func TestGetValidatorOptInStatuses(t *testing.T) {
	m := makeTestMonitor()

	// simulate RPC failure
	m.optChecker = &fakeOptIn{statuses: nil, err: errors.New("rpc fail")}
	duties := []api.ProposerDutyInfo{{PubKey: "pk1"}, {PubKey: "pk2"}}
	res := m.getValidatorOptInStatuses(context.Background(), duties)
	assert.Empty(t, res)

	// success path
	statuses := []validatorrouter.IValidatorOptInRouterOptInStatus{{IsVanillaOptedIn: true}, {IsAvsOptedIn: true}}
	m.optChecker = &fakeOptIn{statuses: statuses, err: nil}
	res = m.getValidatorOptInStatuses(context.Background(), duties)
	assert.Len(t, res, 2)
	assert.True(t, res["pk1"].IsVanillaOptedIn)
	assert.True(t, res["pk2"].IsAvsOptedIn)
}

func TestGetBlocksInfoForDuties(t *testing.T) {
	m := makeTestMonitor()
	m.beacon = &fakeBeaconSlots{}
	duties := []api.ProposerDutyInfo{{Slot: 1}, {Slot: 2}}

	res := m.getBlocksInfoForDuties(context.Background(), duties)
	assert.Equal(t, "10", res[1])
	_, ok := res[2]
	assert.False(t, ok) // slot 2 returned error, should be absent
}

func TestProcessDutySkipsUnopted(t *testing.T) {
	m := makeTestMonitor()
	duty := api.ProposerDutyInfo{PubKey: "pk", Slot: 1, ValidatorIndex: 7}
	notifier := &fakeNotifier{}
	m.notifier = notifier

	m.processDuty(context.Background(), duty, map[string]validatorrouter.IValidatorOptInRouterOptInStatus{}, map[uint64]string{1: "10"})
	assert.False(t, notifier.called)
}

func TestProcessDuty_NotifiesOnOptedIn(t *testing.T) {
	m := makeTestMonitor()
	duty := api.ProposerDutyInfo{PubKey: "pk", Slot: 1, ValidatorIndex: 7}
	notifier := &fakeNotifier{}
	m.notifier = notifier
	// disable dashboard to avoid nil pointer in log
	m.dashboard = nil

	m.processDuty(
		context.Background(),
		duty,
		map[string]validatorrouter.IValidatorOptInRouterOptInStatus{"pk": {IsVanillaOptedIn: true}},
		map[uint64]string{1: "100"},
	)
	assert.True(t, notifier.called)
}

func TestProcessBlockDataAndSave(t *testing.T) {
	m := makeTestMonitor()
	duty := api.ProposerDutyInfo{PubKey: "pk", Slot: 3, ValidatorIndex: 42}

	m.relay = &fakeRelay{
		results: map[string]api.RelayResult{
			"r1": {Response: []api.BidTrace{{ProposerPubkey: "pk", Value: "123", NumTx: "1"}}},
		},
	}
	m.dashboard = &fakeDashboard{
		resp: &api.DashboardResponse{
			Winner:                 "pk",
			TotalOpenedCommitments: 5,
			TotalRewards:           100,
		},
		err: nil,
	}
	notifier := &fakeNotifier{}
	m.notifier = notifier
	fdb := &fakeDB{}
	m.db = fdb

	m.processBlockData(context.Background(), 77, duty)

	assert.True(t, notifier.called)
	assert.Len(t, fdb.saved, 1)
	rec := fdb.saved[0]
	assert.Equal(t, uint64(3), rec.Slot)
	assert.Equal(t, uint64(77), rec.BlockNumber)
	assert.Equal(t, big.NewInt(123), rec.MEVReward)
	assert.Equal(t, []string{"r1"}, rec.RelaysWithData)
	assert.Equal(t, "pk", rec.Winner)
	assert.Equal(t, 5, rec.TotalCommitments)
	assert.Equal(t, 100, rec.TotalRewards)
}

func TestFetchDutiesForEpoch(t *testing.T) {
	m := makeTestMonitor()
	m.beacon = &fakeBeacon{resp: &api.ProposerDutiesResponse{}, err: nil}

	duties, err := m.fetchDutiesForEpoch(context.Background(), 10)
	assert.NoError(t, err)
	assert.Empty(t, duties)
}

func TestFetchDutiesForEpoch_Error(t *testing.T) {
	m := makeTestMonitor()
	m.beacon = &fakeBeacon{resp: nil, err: errors.New("rpc failure")}

	duties, err := m.fetchDutiesForEpoch(context.Background(), 5)
	assert.Error(t, err)
	assert.Nil(t, duties)
}

func TestFetchAndProcessDuties_Caches(t *testing.T) {
	m := makeTestMonitor()
	resp := &api.ProposerDutiesResponse{}
	m.beacon = &fakeBeacon{resp: resp, err: nil}
	m.calculator = &fakeCalc{curEpoch: 2, toFetch: []uint64{1, 2}}

	m.fetchAndProcessDuties(context.Background())
	assert.Contains(t, m.dutiesCache, uint64(1))
	assert.Contains(t, m.dutiesCache, uint64(2))

	// second call should use cache (nothing panics, cache unchanged)
	m.fetchAndProcessDuties(context.Background())
}

func TestFetchBlockInfoFromDashboard_Error(t *testing.T) {
	m := makeTestMonitor()
	m.dashboard = &fakeDashboard{resp: nil, err: errors.New("down")}

	info := m.fetchBlockInfoFromDashboard(context.Background(), 99)
	assert.Nil(t, info)
}

func TestCleanupCaches_Threshold(t *testing.T) {
	m := makeTestMonitor()
	// fill processedBlocks > 501 entries
	for i := range 600 {
		m.processedBlocks[uint64(i)] = time.Now()
	}
	// add one fresh epoch so dutiesCache isn't empty
	m.dutiesCache[5] = cachedDuties{cachedAt: time.Now()}

	m.cleanupCaches()

	assert.Less(t, len(m.processedBlocks), 600) // cleaned
}

// Test fetching and saving commitments
func TestFetchAndSaveCommitments(t *testing.T) {
	m := makeTestMonitor()

	// Create mock commitments
	mockCommitments := []api.CommitmentData{
		{
			CommitmentIndex:     [32]byte{1, 2, 3},
			Bidder:              "0xbidder1",
			Committer:           "0xcommitter1",
			BidAmt:              "1000000000",
			SlashAmt:            "500000000",
			BlockNumber:         42,
			DecayStartTimeStamp: 100,
			DecayEndTimeStamp:   200,
			TxnHash:             "0xtxhash1",
			CommitmentDigest:    [32]byte{4, 5, 6},
		},
		{
			CommitmentIndex:     [32]byte{7, 8, 9},
			Bidder:              "0xbidder2",
			Committer:           "0xcommitter2",
			BidAmt:              "2000000000",
			SlashAmt:            "1000000000",
			BlockNumber:         42,
			DecayStartTimeStamp: 100,
			DecayEndTimeStamp:   200,
			TxnHash:             "0xtxhash2",
			CommitmentDigest:    [32]byte{10, 11, 12},
		},
	}

	// Set up the dashboard mock
	dashboard := &fakeDashboard{
		commitments: mockCommitments,
		commitErr:   nil,
	}
	m.dashboard = dashboard

	// Set up the DB mock
	db := &fakeDB{}
	m.db = db

	// Call the method under test
	m.fetchAndSaveCommitments(context.Background(), 42)

	// Verify results
	assert.Len(t, db.commitments, 2)
	assert.Equal(t, "0xbidder1", db.commitments[0].Bidder)
	assert.Equal(t, "0xcommitter1", db.commitments[0].Committer)
	assert.Equal(t, big.NewInt(1000000000), db.commitments[0].BidAmount)
	assert.Equal(t, "0xbidder2", db.commitments[1].Bidder)
}

// Test handling dashboard API errors
func TestFetchAndSaveCommitments_DashboardError(t *testing.T) {
	m := makeTestMonitor()

	// Set up the dashboard mock with error
	dashboard := &fakeDashboard{
		commitments: nil,
		commitErr:   errors.New("dashboard API error"),
	}
	m.dashboard = dashboard

	// Set up the DB mock
	db := &fakeDB{}
	m.db = db

	// Call the method under test
	m.fetchAndSaveCommitments(context.Background(), 42)

	// Verify nothing saved
	assert.Len(t, db.commitments, 0)
}

// Test handling empty commitments response
func TestFetchAndSaveCommitments_EmptyResponse(t *testing.T) {
	m := makeTestMonitor()

	// Set up the dashboard mock with empty response
	dashboard := &fakeDashboard{
		commitments: []api.CommitmentData{},
		commitErr:   nil,
	}
	m.dashboard = dashboard

	// Set up the DB mock
	db := &fakeDB{}
	m.db = db

	// Call the method under test
	m.fetchAndSaveCommitments(context.Background(), 42)

	// Verify nothing saved
	assert.Len(t, db.commitments, 0)
}

// Test handling database errors
func TestFetchAndSaveCommitments_DBError(t *testing.T) {
	m := makeTestMonitor()

	// Create mock commitments
	mockCommitments := []api.CommitmentData{
		{
			CommitmentIndex:     [32]byte{1, 2, 3},
			Bidder:              "0xbidder1",
			Committer:           "0xcommitter1",
			BidAmt:              "1000000000",
			SlashAmt:            "500000000",
			BlockNumber:         42,
			DecayStartTimeStamp: 100,
			DecayEndTimeStamp:   200,
			TxnHash:             "0xtxhash1",
			CommitmentDigest:    [32]byte{4, 5, 6},
		},
	}

	// Set up the dashboard mock
	dashboard := &fakeDashboard{
		commitments: mockCommitments,
		commitErr:   nil,
	}
	m.dashboard = dashboard

	// Set up the DB mock with error
	db := &fakeDB{
		commitErr: errors.New("database error"),
	}
	m.db = db

	// Call the method under test
	m.fetchAndSaveCommitments(context.Background(), 42)

	// Verify nothing saved due to error
	assert.Len(t, db.commitments, 0)
}

// Test integration with processBlockData
func TestProcessBlockData_FetchesCommitments(t *testing.T) {
	m := makeTestMonitor()
	duty := api.ProposerDutyInfo{PubKey: "pk", Slot: 3, ValidatorIndex: 42}

	// Setup relay mock
	m.relay = &fakeRelay{
		results: map[string]api.RelayResult{
			"r1": {Response: []api.BidTrace{{ProposerPubkey: "pk", Value: "123", NumTx: "1"}}},
		},
	}

	// Create mock commitments
	mockCommitments := []api.CommitmentData{
		{
			CommitmentIndex:     [32]byte{1, 2, 3},
			Bidder:              "0xbidder1",
			Committer:           "0xcommitter1",
			BidAmt:              "1000000000",
			SlashAmt:            "500000000",
			BlockNumber:         77,
			DecayStartTimeStamp: 100,
			DecayEndTimeStamp:   200,
			TxnHash:             "0xtxhash1",
			CommitmentDigest:    [32]byte{4, 5, 6},
		},
	}

	// Setup dashboard mock
	dashboard := &fakeDashboard{
		resp: &api.DashboardResponse{
			Winner:                 "pk",
			TotalOpenedCommitments: 5,
			TotalRewards:           100,
		},
		err:         nil,
		commitments: mockCommitments,
		commitErr:   nil,
	}
	m.dashboard = dashboard

	// Setup notification mock
	notifier := &fakeNotifier{}
	m.notifier = notifier

	// Setup DB mock
	fdb := &fakeDB{}
	m.db = fdb

	// Call the method under test
	m.processBlockData(context.Background(), 77, duty)

	// Verify results
	assert.True(t, notifier.called)
	assert.Len(t, fdb.saved, 1)
	assert.Len(t, fdb.commitments, 1)
	assert.Equal(t, "0xbidder1", fdb.commitments[0].Bidder)
	assert.Equal(t, "0xcommitter1", fdb.commitments[0].Committer)
	assert.Equal(t, big.NewInt(1000000000), fdb.commitments[0].BidAmount)
}
