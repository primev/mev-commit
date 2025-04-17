package notifier

import (
	"github.com/ethereum/go-ethereum/core/types"
)

func (s *FullNotifier) HandleHeader(header *types.Header) error {
	return s.handleHeader(header)
}
