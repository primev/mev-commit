package member

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"log/slog"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	gtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/cl/ethclient"
	"github.com/primev/mev-commit/cl/pb/pb"
	"github.com/primev/mev-commit/cl/redisapp/blockbuilder"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MemberClient struct {
	clientID     string
	streamerAddr string
	conn         *grpc.ClientConn
	client       pb.PayloadStreamerClient
	logger       *slog.Logger
	engineCl     EngineClient
	bb           BlockBuilder
}

type EngineClient interface {
	NewPayloadV3(ctx context.Context, params engine.ExecutableData, versionedHashes []common.Hash, beaconRoot *common.Hash) (engine.PayloadStatusV1, error)
	ForkchoiceUpdatedV3(ctx context.Context, update engine.ForkchoiceStateV1, payloadAttributes *engine.PayloadAttributes) (engine.ForkChoiceResponse, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*gtypes.Header, error)
}

type BlockBuilder interface {
	FinalizeBlock(ctx context.Context, payloadIDStr, executionPayloadStr, msgID string) error
}

func NewMemberClient(clientID, streamerAddr, ecURL, jwtSecret string, logger *slog.Logger) (*MemberClient, error) {
	conn, err := grpc.NewClient(streamerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewPayloadStreamerClient(conn)

	bytes, err := hex.DecodeString(jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("error decoding JWT secret: %v", err)
	}

	engineCL, err := ethclient.NewAuthClient(context.Background(), ecURL, bytes)
	if err != nil {
		return nil, fmt.Errorf("error creating engine client: %v", err)
	}

	bb := blockbuilder.NewMemberBlockBuilder(engineCL, logger)

	return &MemberClient{
		clientID:     clientID,
		streamerAddr: streamerAddr,
		conn:         conn,
		client:       client,
		engineCl:     engineCL,
		logger:       logger,
		bb:           bb,
	}, nil
}

func (mc *MemberClient) Run(ctx context.Context) error {
	stream, err := mc.client.Subscribe(ctx)
	if err != nil {
		return err
	}

	err = stream.Send(&pb.ClientMessage{
		Message: &pb.ClientMessage_SubscribeRequest{
			SubscribeRequest: &pb.SubscribeRequest{
				ClientId: mc.clientID,
			},
		},
	})
	if err != nil {
		mc.logger.Error("Failed to send SubscribeRequest", "error", err)
		return err
	}

	mc.logger.Info("Member client started", "clientID", mc.clientID)

	for {
		select {
		case <-ctx.Done():
			mc.logger.Info("Member client context done", "clientID", mc.clientID)
			return nil
		default:
			msg, err := stream.Recv()
			if err != nil {
				if errors.Is(err, context.Canceled) {
					mc.logger.Info("Member client context canceled", "clientID", mc.clientID)
					return nil
				}
				mc.logger.Error("Error receiving message", "error", err)
				continue
			}
			err = mc.bb.FinalizeBlock(ctx, msg.PayloadId, msg.ExecutionPayload, msg.MessageId)
			if err != nil {
				mc.logger.Error("Error processing payload", "error", err)
				continue
			}

			err = stream.Send(&pb.ClientMessage{
				Message: &pb.ClientMessage_AckPayload{
					AckPayload: &pb.AckPayloadRequest{
						ClientId:  mc.clientID,
						PayloadId: msg.PayloadId,
						MessageId: msg.MessageId,
					},
				},
			})
			if err != nil {
				mc.logger.Error("Failed to send acknowledgment", "error", err)
				continue
			}
			mc.logger.Info("Acknowledged message", "payloadID", msg.PayloadId)
		}
	}
}
