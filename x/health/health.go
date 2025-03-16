package health

import (
	"fmt"

	"google.golang.org/grpc/connectivity"
)

type HealthCheck interface {
	// Check checks the health of the service.
	Check() error
}

type Health interface {
	// Register registers a health check.
	Register(HealthCheck)
	// Health returns the health of the service.
	Health() error
}

type HealthCheckFunc func() error

func (h HealthCheckFunc) Check() error {
	return h()
}

// CloseChannelHealthCheck returns a health check that checks if the channel is closed.
func CloseChannelHealthCheck(desc string, ch <-chan struct{}) HealthCheck {
	return HealthCheckFunc(func() error {
		select {
		case <-ch:
			return fmt.Errorf("%s: closed", desc)
		default:
			return nil
		}
	})
}

// GrpcClientConn is an interface that represents the methods we need from grpc.ClientConn.
type GrpcClientConn interface {
	GetState() connectivity.State
}

// GrpcGatewayHealthCheck returns a health check that checks the state of the gRPC connection.
func GrpcGatewayHealthCheck(conn GrpcClientConn) HealthCheck {
	return HealthCheckFunc(func() error {
		if conn.GetState() == connectivity.TransientFailure {
			return fmt.Errorf("grpc gateway: %s", conn.GetState())
		}
		return nil
	})
}

func New() Health {
	return &health{}
}

type health struct {
	checks []HealthCheck
}

func (h *health) Register(check HealthCheck) {
	h.checks = append(h.checks, check)
}

func (h *health) Health() error {
	for _, check := range h.checks {
		if err := check.Check(); err != nil {
			return err
		}
	}
	return nil
}
