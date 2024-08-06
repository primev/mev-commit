package validatorapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	validatorapiv1 "github.com/primev/mev-commit/p2p/gen/go/validatorapi/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ValidatorRouterContract interface {
	AreValidatorsOptedIn(valBLSPubKeys [][]byte) ([]bool, error)
}

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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.apiURL+"/eth/v1/beacon/states/head/finality_checkpoints", nil)
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
		return nil, status.Errorf(codes.Internal, "unexpected status code: %v", resp.StatusCode)
	}

	var dutiesResp ProposerDutiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&dutiesResp); err != nil {
		s.logger.Error("decoding response", "error", err)
		return nil, status.Errorf(codes.Internal, "decoding response: %v", err)
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

	areValidatorsOptedIn, err := s.validatorRouter.AreValidatorsOptedIn(validatorsKeys)
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
			BLSKey:    duty.Pubkey,
			IsOptedIn: areValidatorsOptedIn[i],
		}
	}

	return validators, nil
}