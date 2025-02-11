package stakemanager_test

import (
	"context"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	"github.com/primev/mev-commit/p2p/pkg/notifications"
	"github.com/primev/mev-commit/p2p/pkg/stakemanager"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/util"
)

type mockProviderRegistry struct {
	providerStakes map[common.Address]*big.Int
	minStake       *big.Int
}

func (m *mockProviderRegistry) GetProviderStake(_ *bind.CallOpts, addr common.Address) (*big.Int, error) {
	stake, found := m.providerStakes[addr]
	if !found {
		return big.NewInt(0), nil
	}

	return stake, nil
}

func (m *mockProviderRegistry) MinStake(_ *bind.CallOpts) (*big.Int, error) {
	return new(big.Int).Set(m.minStake), nil
}

type mockNotifier struct {
	evt chan *notifications.Notification
}

func (m *mockNotifier) Notify(n *notifications.Notification) {
	m.evt <- n
}

func TestStakeManager(t *testing.T) {
	t.Parallel()

	providerStakes := map[common.Address]*big.Int{
		common.HexToAddress("0x456"): big.NewInt(100),
		common.HexToAddress("0x789"): big.NewInt(200),
		common.HexToAddress("0xabc"): big.NewInt(500),
	}

	owner := common.HexToAddress("0x123")
	providerRegistry := &mockProviderRegistry{
		providerStakes: providerStakes,
		minStake:       big.NewInt(10),
	}

	notifier := &mockNotifier{
		evt: make(chan *notifications.Notification, 10),
	}

	prABI, err := abi.JSON(strings.NewReader(providerregistry.ProviderregistryABI))
	if err != nil {
		t.Fatalf("failed to parse provider registry ABI: %v", err)
	}
	evtMgr := events.NewListener(
		util.NewTestLogger(os.Stdout),
		&prABI,
	)

	stakeMgr, err := stakemanager.NewStakeManager(
		owner,
		evtMgr,
		providerRegistry,
		notifier,
		util.NewTestLogger(os.Stdout),
	)
	if err != nil {
		t.Fatalf("failed to create stake manager: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := stakeMgr.Start(ctx)

	minStake := stakeMgr.MinStake()
	if minStake.Cmp(big.NewInt(10)) != 0 {
		t.Errorf("unexpected min stake: %v", minStake)
	}

	// Simulates getting the value from the contract
	stake, err := stakeMgr.GetStake(common.HexToAddress("0x456"))
	if err != nil {
		t.Fatalf("failed to get stake: %v", err)
	}

	if stake.Cmp(big.NewInt(100)) != 0 {
		t.Errorf("unexpected stake: %v", stake)
	}

	if !stakeMgr.CheckProviderRegistered(common.HexToAddress("0x456")) {
		t.Errorf("provider should be registered")
	}

	publishMinStakeUpdated(
		evtMgr,
		&prABI,
		providerregistry.ProviderregistryMinStakeUpdated{
			NewMinStake: big.NewInt(20),
		},
	)

	start := time.Now()
	for stakeMgr.MinStake().Cmp(big.NewInt(20)) != 0 {
		if time.Since(start) > 5*time.Second {
			t.Fatalf("timed out waiting for min stake to update")
		}
		time.Sleep(100 * time.Millisecond)
	}

	for addr, stake := range providerStakes {
		publishProviderRegistered(
			evtMgr,
			&prABI,
			providerregistry.ProviderregistryProviderRegistered{
				Provider:     addr,
				StakedAmount: stake,
			},
		)
	}

	count := 0
	for e := range notifier.evt {
		if e.Topic() == notifications.TopicProviderRegistered {
			count++
		}
		// We will get the event only for unknown providers
		if count == len(providerStakes)-1 {
			break
		}
	}

	for addr, stake := range providerStakes {
		providerStakes[addr] = new(big.Int).Sub(stake, big.NewInt(10))
		publishFundsSlashed(
			evtMgr,
			&prABI,
			providerregistry.ProviderregistryFundsSlashed{
				Provider: addr,
				Amount:   big.NewInt(10),
			},
		)
	}

	count = 0
	for e := range notifier.evt {
		if e.Topic() == notifications.TopicProviderSlashed {
			count++
		}
		addr := common.HexToAddress(e.Value()["provider"].(string))
		stake, err := stakeMgr.GetStake(addr)
		if err != nil {
			t.Fatalf("failed to get stake: %v", err)
		}
		if stake.Cmp(providerStakes[addr]) != 0 {
			t.Errorf("unexpected stake: %v", stake)
		}
		if count == len(providerStakes) {
			break
		}
	}

	for addr, stake := range providerStakes {
		providerStakes[addr] = new(big.Int).Add(stake, big.NewInt(10))
		publishFundsDeposited(
			evtMgr,
			&prABI,
			providerregistry.ProviderregistryFundsDeposited{
				Provider: addr,
				Amount:   big.NewInt(10),
			},
		)
	}

	count = 0
	for e := range notifier.evt {
		if e.Topic() == notifications.TopicProviderDeposit {
			count++
		}
		addr := common.HexToAddress(e.Value()["provider"].(string))
		stake, err := stakeMgr.GetStake(addr)
		if err != nil {
			t.Fatalf("failed to get stake: %v", err)
		}
		if stake.Cmp(providerStakes[addr]) != 0 {
			t.Errorf("unexpected stake: %v", stake)
		}
		if count == len(providerStakes) {
			break
		}
	}

	for addr, _ := range providerStakes {
		if !stakeMgr.CheckProviderRegistered(addr) {
			t.Errorf("provider should be registered")
		}
	}

	publishUnstake(
		evtMgr,
		&prABI,
		providerregistry.ProviderregistryUnstake{
			Provider:  common.HexToAddress("0x456"),
			Timestamp: big.NewInt(0),
		},
	)

	e := <-notifier.evt
	if e.Topic() != notifications.TopicProviderDeregistered {
		t.Errorf("unexpected notification: %v", e)
	}

	if stakeMgr.CheckProviderRegistered(common.HexToAddress("0x456")) {
		t.Errorf("provider should not be registered")
	}

	cancel()
	<-done
}

func publishProviderRegistered(
	evtMgr events.EventManager,
	prABI *abi.ABI,
	ev providerregistry.ProviderregistryProviderRegistered,
) error {
	event := prABI.Events["ProviderRegistered"]
	buf, err := event.Inputs.NonIndexed().Pack(
		ev.StakedAmount,
	)
	if err != nil {
		return err
	}

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			event.ID,                            // The first topic is the hash of the event signature
			common.HexToHash(ev.Provider.Hex()), // The next topics are the indexed event parameters
		},
		// Non-indexed parameters are stored in the Data field
		Data: buf,
	}

	evtMgr.PublishLogEvent(context.Background(), testLog)
	return nil
}

func publishFundsSlashed(
	evMgr events.EventManager,
	prABI *abi.ABI,
	ev providerregistry.ProviderregistryFundsSlashed,
) error {
	event := prABI.Events["FundsSlashed"]
	buf, err := event.Inputs.NonIndexed().Pack(
		ev.Amount,
	)
	if err != nil {
		return err
	}

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			event.ID,                            // The first topic is the hash of the event signature
			common.HexToHash(ev.Provider.Hex()), // The next topics are the indexed event parameters
		},
		// Non-indexed parameters are stored in the Data field
		Data: buf,
	}

	evMgr.PublishLogEvent(context.Background(), testLog)
	return nil
}

func publishFundsDeposited(
	evMgr events.EventManager,
	prABI *abi.ABI,
	ev providerregistry.ProviderregistryFundsDeposited,
) error {
	event := prABI.Events["FundsDeposited"]
	buf, err := event.Inputs.NonIndexed().Pack(
		ev.Amount,
	)
	if err != nil {
		return err
	}

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			event.ID,                            // The first topic is the hash of the event signature
			common.HexToHash(ev.Provider.Hex()), // The next topics are the indexed event parameters
		},
		// Non-indexed parameters are stored in the Data field
		Data: buf,
	}

	evMgr.PublishLogEvent(context.Background(), testLog)
	return nil
}

func publishUnstake(
	evMgr events.EventManager,
	prABI *abi.ABI,
	ev providerregistry.ProviderregistryUnstake,
) error {
	event := prABI.Events["Unstake"]
	buf, err := event.Inputs.NonIndexed().Pack(
		ev.Timestamp,
	)
	if err != nil {
		return err
	}

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			event.ID,                            // The first topic is the hash of the event signature
			common.HexToHash(ev.Provider.Hex()), // The next topics are the indexed event parameters
		},
		// Non-indexed parameters are stored in the Data field
		Data: buf,
	}

	evMgr.PublishLogEvent(context.Background(), testLog)
	return nil
}

func publishMinStakeUpdated(
	evMgr events.EventManager,
	prABI *abi.ABI,
	ev providerregistry.ProviderregistryMinStakeUpdated,
) error {
	event := prABI.Events["MinStakeUpdated"]

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			event.ID, // The first topic is the hash of the event signature
			common.BigToHash(ev.NewMinStake),
		},
		// Non-indexed parameters are stored in the Data field
		Data: nil,
	}

	evMgr.PublishLogEvent(context.Background(), testLog)
	return nil
}
