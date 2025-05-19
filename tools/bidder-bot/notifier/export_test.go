package notifier

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
)

func (s *FullNotifier) HandleHeader(ctx context.Context, header *types.Header) error {
	return s.handleHeader(ctx, header)
}

func (s *SelectiveNotifier) HandleMsg(ctx context.Context, msg *notificationsapiv1.Notification) error {
	return s.handleMsg(ctx, msg)
}
