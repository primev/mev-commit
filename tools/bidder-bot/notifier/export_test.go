package notifier

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

func (s *FullNotifier) HandleHeader(ctx context.Context, header *types.Header) error {
	return s.handleHeader(ctx, header)
}
