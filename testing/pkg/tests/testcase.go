package tests

import (
	"context"

	"github.com/primev/mev-commit/testing/pkg/orchestrator"
	"github.com/primev/mev-commit/testing/pkg/tests/connectivity"
	"github.com/primev/mev-commit/testing/pkg/tests/deposit"
	"github.com/primev/mev-commit/testing/pkg/tests/preconf"
	"github.com/primev/mev-commit/testing/pkg/tests/staking"
)

type TestCase func(ctx context.Context, cluster orchestrator.Orchestrator, args any) error

var TestCases = map[string]TestCase{
	"staking":           staking.Run,
	"connectivity":      connectivity.Run,
	"autodeposit":       deposit.RunAutoDeposit,
	"preconf":           preconf.RunPreconf,
	"cancelAutodeposit": deposit.RunCancelAutoDeposit,
}
