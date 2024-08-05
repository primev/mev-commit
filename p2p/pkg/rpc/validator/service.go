package validatorapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	validatorapiv1 "github.com/primev/mev-commit/p2p/gen/go/validatorapi/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ValidatorRouterContract interface {
	AreValidatorsOptedIn(valBLSPubKeys [][]byte) ([]bool, error)
}

type OptsGetter func(context.Context) (*bind.CallOpts, error)

type Service struct {
	validatorapiv1.UnimplementedValidatorServer
	apiURL          string
	validatorRouter ValidatorRouterContract
	logger          *slog.Logger
	metrics         *metrics
}

func NewService(apiURL string, validatorRouter ValidatorRouterContract, logger *slog.Logger) *Service {
	return &Service{
		apiURL:          apiURL,
		validatorRouter: validatorRouter,
		logger:          logger,
		metrics:         newMetrics(),
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

func (s *Service) GetEpoch(
	ctx context.Context,
	_ *validatorapiv1.EmptyMessage,
) (*validatorapiv1.EpochResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.apiURL+"/eth/v1/beacon/states/head/finality_checkpoints", nil)
	if err != nil {
		s.logger.Error("creating request", "error", err)
		return nil, status.Errorf(codes.Internal, "creating request: %v", err)
	}

	req.Header.Set("accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("making request", "error", err)
		return nil, status.Errorf(codes.Internal, "making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("unexpected status code", "status", resp.StatusCode)
		return nil, status.Errorf(codes.Internal, "unexpected status code: %v", resp.StatusCode)
	}

	var checkpointsResp FinalityCheckpointsResponse
	if err := json.NewDecoder(resp.Body).Decode(&checkpointsResp); err != nil {
		s.logger.Error("decoding response", "error", err)
		return nil, status.Errorf(codes.Internal, "decoding response: %v", err)
	}

	currentEpoch := parseEpoch(checkpointsResp.Data.CurrentJustified.Epoch)
	previousJustifiedEpoch := parseEpoch(checkpointsResp.Data.PreviousJustified.Epoch)
	finalizedEpoch := parseEpoch(checkpointsResp.Data.Finalized.Epoch)

	s.logger.Info("fetched epoch data", "current_epoch", currentEpoch, "current_justified_epoch", currentEpoch, "previous_justified_epoch", previousJustifiedEpoch, "finalized_epoch", finalizedEpoch)
	s.metrics.FetchedEpochDataCount.Inc()

	return &validatorapiv1.EpochResponse{
		CurrentEpoch:           currentEpoch + 1,
		CurrentJustifiedEpoch:  currentEpoch,
		PreviousJustifiedEpoch: previousJustifiedEpoch,
		FinalizedEpoch:         finalizedEpoch,
	}, nil
}

func parseEpoch(epochStr string) uint64 {
	var epoch uint64
	fmt.Sscanf(epochStr, "%d", &epoch)
	return epoch
}

func (s *Service) GetValidators(
	ctx context.Context,
	req *validatorapiv1.GetValidatorsRequest,
) (*validatorapiv1.GetValidatorsResponse, error) {
	epoch := req.Epoch

	// If epoch is zero, fetch the current epoch
	if epoch == 0 {
		epochResponse, err := s.GetEpoch(ctx, &validatorapiv1.EmptyMessage{})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "fetching current epoch: %v", err)
		}
		epoch = epochResponse.CurrentEpoch
	}

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
		return nil, status.Errorf(codes.Internal, "unexpected status code: %v", resp.StatusCode)
	}

	var dutiesResp ProposerDutiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&dutiesResp); err != nil {
		s.logger.Error("decoding response", "error", err)
		return nil, status.Errorf(codes.Internal, "decoding response: %v", err)
	}

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
	areValidatorsOptedIn, err := s.validatorRouter.AreValidatorsOptedIn(validatorsKeys[:2])
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
		validators[slot] = &validatorapiv1.SlotInfo{
			BLSKey:   duty.Pubkey,
			IsActive: areValidatorsOptedIn[i],
		}
	}

	s.logger.Info("fetched validators for epoch", "epoch", epoch, "validators", validators)
	s.metrics.FetchedValidatorsCount.Inc()

	return &validatorapiv1.GetValidatorsResponse{
		Items: validators,
	}, nil
}
