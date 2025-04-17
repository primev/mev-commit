package validatorapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	validatoroptinrouter "github.com/primev/mev-commit/contracts-abi/clients/ValidatorOptInRouter"
	validatorapiv1 "github.com/primev/mev-commit/p2p/gen/go/validatorapi/v1"
	"github.com/primev/mev-commit/p2p/pkg/notifications"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ValidatorRouterContract interface {
	AreValidatorsOptedIn(opts *bind.CallOpts, valBLSPubKeys [][]byte) ([]validatoroptinrouter.IValidatorOptInRouterOptInStatus, error)
}

type Service struct {
	validatorapiv1.UnimplementedValidatorServer
	apiURL               string
	validatorRouter      ValidatorRouterContract
	logger               *slog.Logger
	metrics              *metrics
	optsGetter           func() (*bind.CallOpts, error)
	notifier             notifications.Notifier
	genesisTime          time.Time
	proposerNotifyOffset time.Duration

	slotDuration  time.Duration
	epochSlots    int
	epochDuration time.Duration
}

const (
	SlotDuration  = 12 * time.Second
	EpochSlots    = 32
	EpochDuration = SlotDuration * time.Duration(EpochSlots)
)

func NewService(
	apiURL string,
	validatorRouter ValidatorRouterContract,
	logger *slog.Logger,
	optsGetter func() (*bind.CallOpts, error),
	notifier notifications.Notifier,
	proposerNotifyOffset time.Duration,
) *Service {
	return &Service{
		apiURL:               apiURL,
		validatorRouter:      validatorRouter,
		logger:               logger,
		metrics:              newMetrics(),
		optsGetter:           optsGetter,
		notifier:             notifier,
		proposerNotifyOffset: proposerNotifyOffset,
		slotDuration:         SlotDuration,
		epochSlots:           EpochSlots,
		epochDuration:        EpochDuration,
	}
}

type FinalityCheckpointsResponse struct {
	Data struct {
		PreviousJustified struct {
			Epoch string `json:"epoch"`
		} `json:"previous_justified"`
		CurrentJustified struct {
			Epoch string `json:"epoch"`
		} `json:"current_justified"`
		Finalized struct {
			Epoch string `json:"epoch"`
		} `json:"finalized"`
	} `json:"data"`
}

type ProposerDutiesResponse struct {
	Data []struct {
		Pubkey string `json:"pubkey"`
		Slot   string `json:"slot"`
	} `json:"data"`
}

func (s *Service) GetValidators(
	ctx context.Context,
	req *validatorapiv1.GetValidatorsRequest,
) (*validatorapiv1.GetValidatorsResponse, error) {
	currentEpoch, err := s.fetchCurrentEpoch(ctx, req.Epoch)
	if err != nil {
		return nil, err
	}

	dutiesResp, err := s.fetchProposerDuties(ctx, currentEpoch)
	if err != nil {
		return nil, err
	}

	validators, err := s.processValidators(dutiesResp)
	if err != nil {
		return nil, err
	}

	s.logger.Info("fetched validators for epoch", "current_epoch", currentEpoch, "validators", validators)
	s.metrics.FetchedValidatorsCount.Inc()

	return &validatorapiv1.GetValidatorsResponse{
		Items: validators,
	}, nil
}

func (s *Service) fetchCurrentEpoch(ctx context.Context, epoch uint64) (uint64, error) {
	if epoch != 0 {
		return epoch, nil
	}
	url := fmt.Sprintf("%s/eth/v1/beacon/states/head/finality_checkpoints", s.apiURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		s.logger.Error("creating request", "error", err)
		return 0, status.Errorf(codes.Internal, "creating request: %v", err)
	}

	req.Header.Set("accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("making request", "error", err)
		return 0, status.Errorf(codes.Internal, "making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("unexpected status code", "status", resp.StatusCode)
		return 0, status.Errorf(codes.Internal, "unexpected status code: %v", resp.StatusCode)
	}

	var checkpointsResp FinalityCheckpointsResponse
	if err := json.NewDecoder(resp.Body).Decode(&checkpointsResp); err != nil {
		s.logger.Error("decoding response", "error", err)
		return 0, status.Errorf(codes.Internal, "decoding response: %v", err)
	}

	currentJustifiedEpoch, err := strconv.ParseUint(checkpointsResp.Data.CurrentJustified.Epoch, 10, 64)
	if err != nil {
		return 0, status.Errorf(codes.Internal, "parsing current justified epoch: %v", err)
	}

	return currentJustifiedEpoch + 1, nil
}

func (s *Service) fetchProposerDuties(ctx context.Context, epoch uint64) (*ProposerDutiesResponse, error) {
	url := fmt.Sprintf("%s/eth/v1/validator/duties/proposer/%d", s.apiURL, epoch)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		s.logger.Error("creating request", "error", err)
		return nil, status.Errorf(codes.Internal, "creating request: %v", err)
	}

	httpReq.Header.Set("accept", "application/json")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		s.logger.Error("making request", "error", err)
		return nil, status.Errorf(codes.Internal, "making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("unexpected status code", "status", resp.StatusCode)

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "reading response body: %v", err)
		}

		bodyString := string(bodyBytes)
		if strings.Contains(bodyString, "Proposer duties were requested for a future epoch") {
			return nil, status.Errorf(codes.InvalidArgument, "Proposer duties were requested for a future epoch")
		}

		return nil, status.Errorf(
			codes.Internal,
			"unexpected status code: %v, response: %s", resp.StatusCode, bodyString,
		)
	}
	var dutiesResp ProposerDutiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&dutiesResp); err != nil {
		s.logger.Error("decoding response", "error", err)
		return nil, status.Errorf(codes.Internal, "decoding response: %v", err)
	}

	if len(dutiesResp.Data) == 0 {
		return nil, status.Errorf(codes.Internal, "no proposer duties found for epoch %d", epoch)
	}

	return &dutiesResp, nil
}

func (s *Service) processValidators(dutiesResp *ProposerDutiesResponse) (map[uint64]*validatorapiv1.SlotInfo, error) {
	validators := make(map[uint64]*validatorapiv1.SlotInfo, len(dutiesResp.Data))
	validatorsKeys := make([][]byte, 0, len(dutiesResp.Data))
	for _, duty := range dutiesResp.Data {
		pubkeyBytes, err := hexutil.Decode(duty.Pubkey)
		if err != nil {
			s.logger.Error("decoding pubkey", "error", err)
			continue
		}

		validatorsKeys = append(validatorsKeys, pubkeyBytes)
	}

	if len(validatorsKeys) == 0 {
		return validators, nil
	}

	opts, err := s.optsGetter()
	if err != nil {
		s.logger.Error("getting call opts", "error", err)
		return nil, status.Errorf(codes.Internal, "getting call opts: %v", err)
	}

	areValidatorsOptedIn, err := s.validatorRouter.AreValidatorsOptedIn(opts, validatorsKeys)
	if err != nil {
		s.logger.Error("checking if validators are opted in", "error", err)
		return nil, status.Errorf(codes.Internal, "checking if validators are opted in: %v", err)
	}

	for i, duty := range dutiesResp.Data {
		slot, err := strconv.ParseUint(duty.Slot, 10, 64)
		if err != nil {
			s.logger.Error("parsing slot number", "error", err)
			continue
		}
		isOptedIn := areValidatorsOptedIn[i].IsVanillaOptedIn ||
			areValidatorsOptedIn[i].IsAvsOptedIn ||
			areValidatorsOptedIn[i].IsMiddlewareOptedIn
		validators[slot] = &validatorapiv1.SlotInfo{
			BLSKey:    duty.Pubkey,
			IsOptedIn: isOptedIn,
		}
	}

	return validators, nil
}

func (s *Service) scheduleNotificationForSlot(epoch uint64, slot uint64, info *validatorapiv1.SlotInfo) {
	slotStartTime := s.genesisTime.Add(time.Duration(slot) * s.slotDuration)
	notificationTime := slotStartTime.Add(-s.proposerNotifyOffset)

	delay := time.Until(notificationTime)
	if delay <= 0 {
		s.logger.Warn("notification time already passed for slot", "epoch", epoch, "slot", slot)
		return
	}

	time.AfterFunc(delay, func() {
		notif := notifications.NewNotification(
			notifications.TopicValidatorOptedIn,
			map[string]any{
				"epoch":   epoch,
				"slot":    slot,
				"bls_key": info.BLSKey,
			},
		)
		s.notifier.Notify(notif)
		s.logger.Info(
			"sent notification for opted in validator",
			"epoch", epoch,
			"slot", slot,
			"bls_key", info.BLSKey,
		)
	})
}

func (s *Service) processEpoch(ctx context.Context, epoch uint64, epochTime time.Time) {
	s.logger.Info("processing epoch", "epoch", epoch)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.processFirstSlotOfNextEpoch(ctx, epoch+1, epochTime.Add(s.epochDuration))
		s.logger.Debug("processed first slot of next epoch", "epoch", epoch+1)
	}()

	dutiesResp, err := s.fetchProposerDuties(ctx, epoch)
	if err != nil {
		s.logger.Error("failed to fetch proposer duties", "epoch", epoch, "error", err)
		wg.Wait()
		return
	}

	validators, err := s.processValidators(dutiesResp)
	if err != nil {
		s.logger.Error("failed to process validators", "epoch", epoch, "error", err)
		wg.Wait()
		return
	}

	firstSlot := s.getFirstSlot(validators)

	optedInSlots := make([]any, 0)
	for slot, info := range validators {
		if info.IsOptedIn {
			optedInSlots = append(optedInSlots, map[string]any{
				"slot":     slot,
				"bls_key":  info.BLSKey,
				"opted_in": info.IsOptedIn,
			})
		}
		if slot != firstSlot {
			s.scheduleNotificationForSlot(epoch, slot, info)
		}
	}

	// Send the notification even in case of no slots
	notif := notifications.NewNotification(
		notifications.TopicEpochValidatorsOptedIn,
		map[string]any{
			"epoch":            epoch,
			"epoch_start_time": epochTime.Unix(),
			"slots":            optedInSlots,
		},
	)
	s.notifier.Notify(notif)
	s.logger.Info("sent notification for epoch with opted in validators",
		"epoch", epoch,
		"opted_in_slot_count", len(optedInSlots),
		"first_slot_in_epoch", firstSlot,
	)

	wg.Wait()
}

func (s *Service) processFirstSlotOfNextEpoch(ctx context.Context, nextEpoch uint64, nextEpochTime time.Time) {
	nextDutiesResp, err := s.fetchProposerDuties(ctx, nextEpoch)
	if err != nil {
		s.logger.Error("failed to fetch proposer duties", "epoch", nextEpoch, "error", err)
		return
	}
	nextValidators, err := s.processValidators(nextDutiesResp)
	if err != nil {
		s.logger.Error("failed to process validators", "epoch", nextEpoch, "error", err)
		return
	}

	firstSlot := s.getFirstSlot(nextValidators)
	info := nextValidators[firstSlot]
	if info.IsOptedIn {
		s.scheduleNotificationForSlot(nextEpoch, firstSlot, info)
	}
}

// Start starts a background job that fetches and processes an epoch every 384 seconds.
// (384 seconds is the duration of an epoch)
func (s *Service) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	genesisTime, err := s.fetchGenesisTime(ctx)
	if err != nil {
		s.logger.Error("failed to fetch genesis time", "error", err)
		close(doneChan)
		return doneChan
	}
	s.genesisTime = genesisTime
	s.logger.Info("initialized genesis time", "genesis_time", s.genesisTime)

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		currentEpoch := uint64(time.Since(s.genesisTime) / s.epochDuration)

		s.processEpoch(egCtx, currentEpoch, s.genesisTime.Add(time.Duration(currentEpoch)*s.epochDuration))
		for {
			nextEpochStart := s.genesisTime.Add(time.Duration(currentEpoch+1) * s.epochDuration)
			delay := max(time.Until(nextEpochStart), 0)

			timer := time.NewTimer(delay)
			select {
			case <-egCtx.Done():
				if !timer.Stop() {
					<-timer.C
				}
				s.logger.Info("epoch cron job stopped")
				return nil
			case <-timer.C:
				currentEpoch++
				s.logger.Info("processing new epoch", "epoch", currentEpoch)
				s.processEpoch(egCtx, currentEpoch, nextEpochStart)
			}
		}
	})

	go func() {
		defer close(doneChan)
		if err := eg.Wait(); err != nil {
			s.logger.Error("error in epoch cron job", "error", err)
		}
	}()

	return doneChan
}

func (s *Service) fetchGenesisTime(ctx context.Context) (time.Time, error) {
	url := fmt.Sprintf("%s/eth/v1/beacon/genesis", s.apiURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		s.logger.Error("creating genesis request", "error", err)
		return time.Time{}, err
	}
	req.Header.Set("accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("making genesis request", "error", err)
		return time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("unexpected status code for genesis", "status", resp.StatusCode)
		return time.Time{}, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	var genesisResp struct {
		Data struct {
			GenesisTime string `json:"genesis_time"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&genesisResp); err != nil {
		s.logger.Error("decoding genesis response", "error", err)
		return time.Time{}, err
	}
	genesisTimeInt, err := strconv.ParseInt(genesisResp.Data.GenesisTime, 10, 64)
	if err != nil {
		s.logger.Error("parsing genesis time", "error", err)
		return time.Time{}, err
	}
	return time.Unix(genesisTimeInt, 0), nil
}

func (s *Service) getFirstSlot(validators map[uint64]*validatorapiv1.SlotInfo) uint64 {
	firstSlot := uint64(math.MaxUint64)
	for slot := range validators {
		if slot < firstSlot {
			firstSlot = slot
		}
	}
	return firstSlot
}
