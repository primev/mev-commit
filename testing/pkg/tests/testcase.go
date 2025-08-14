package tests

import (
	"context"

	"github.com/primev/mev-commit/testing/pkg/orchestrator"
	"github.com/primev/mev-commit/testing/pkg/tests/bridge"
	"github.com/primev/mev-commit/testing/pkg/tests/connectivity"
	"github.com/primev/mev-commit/testing/pkg/tests/preconf"
	"github.com/primev/mev-commit/testing/pkg/tests/staking"
)

type TestCase func(ctx context.Context, cluster orchestrator.Orchestrator, args any) error

type TestEntry struct {
	Name string
	Run  TestCase
}

var TestCases = []TestEntry{
	{"bridge", bridge.RunBridge},
	{"staking", staking.Run},
	{"staking_add_deposit", staking.RunAddDeposit},
	{"connectivity", connectivity.Run},
	// {"autodeposit", deposit.RunAutoDeposit},
	{"preconf", preconf.RunPreconf},
	// {"cancelAutodeposit", deposit.RunCancelAutoDeposit},
}
