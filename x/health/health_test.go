package health_test

import (
	"fmt"
	"testing"

	"github.com/primev/mev-commit/x/health"
	"google.golang.org/grpc/connectivity"
)

func TestCloseChannelHealthCheck(t *testing.T) {
	t.Parallel()

	t.Run("ClosedChannel", func(t *testing.T) {
		ch := make(chan struct{})
		close(ch)

		check := health.CloseChannelHealthCheck("test", ch)
		err := check.Check()
		if err == nil {
			t.Error("expected error, got nil")
		} else if err.Error() != "test: closed" {
			t.Errorf("expected error message 'test: closed', got '%s'", err.Error())
		}
	})

	t.Run("OpenChannel", func(t *testing.T) {
		ch := make(chan struct{})

		check := health.CloseChannelHealthCheck("test", ch)
		err := check.Check()
		if err != nil {
			t.Errorf("expected no error, got '%s'", err.Error())
		}
	})
}

// MockGrpcClientConn is a mock implementation of the GrpcClientConn interface.
type MockGrpcClientConn struct {
	State connectivity.State
}

func (m *MockGrpcClientConn) GetState() connectivity.State {
	return m.State
}

func TestGrpcGatewayHealthCheck(t *testing.T) {
	t.Parallel()

	t.Run("NotReady", func(t *testing.T) {
		conn := &MockGrpcClientConn{State: connectivity.TransientFailure}

		check := health.GrpcGatewayHealthCheck(conn)
		err := check.Check()
		if err == nil {
			t.Error("expected error, got nil")
		} else if err.Error() != fmt.Sprintf("grpc gateway: %s", conn.GetState()) {
			t.Errorf("expected error message 'grpc gateway: Connecting', got '%s'", err.Error())
		}
	})

	t.Run("Ready", func(t *testing.T) {
		conn := &MockGrpcClientConn{State: connectivity.Ready}

		check := health.GrpcGatewayHealthCheck(conn)
		err := check.Check()
		if err != nil {
			t.Errorf("expected no error, got '%s'", err.Error())
		}
	})
}

func TestHealth(t *testing.T) {
	t.Parallel()

	t.Run("AllChecksPass", func(t *testing.T) {
		h := health.New()

		check1 := func() error { return nil }
		check2 := func() error { return nil }

		h.Register(health.HealthCheckFunc(check1))
		h.Register(health.HealthCheckFunc(check2))

		err := h.Health()
		if err != nil {
			t.Errorf("expected no error, got '%s'", err.Error())
		}
	})

	t.Run("FirstCheckFails", func(t *testing.T) {
		h := health.New()

		check1 := func() error { return nil }
		check2 := func() error { return healthCheckError("check2 failed") }

		h.Register(health.HealthCheckFunc(check1))
		h.Register(health.HealthCheckFunc(check2))

		err := h.Health()
		if err == nil {
			t.Error("expected error, got nil")
		} else if err.Error() != "check2 failed" {
			t.Errorf("expected error message 'check2 failed', got '%s'", err.Error())
		}
	})

	t.Run("SecondCheckFails", func(t *testing.T) {
		h := health.New()

		check1 := func() error { return healthCheckError("check1 failed") }
		check2 := func() error { return nil }

		h.Register(health.HealthCheckFunc(check1))
		h.Register(health.HealthCheckFunc(check2))

		err := h.Health()
		if err == nil {
			t.Error("expected error, got nil")
		} else if err.Error() != "check1 failed" {
			t.Errorf("expected error message 'check1 failed', got '%s'", err.Error())
		}
	})
}

type healthCheckError string

func (e healthCheckError) Error() string {
	return string(e)
}
