package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"time"

	"github.com/primev/mev-commit/tools/validators-monitor/config"
	"github.com/primev/mev-commit/tools/validators-monitor/monitor"
	"github.com/primev/mev-commit/x/health"
)

type Service struct {
	cancel  context.CancelFunc
	closers []io.Closer
	logger  *slog.Logger
}

// New creates a new duty monitor service
func New(cfg *config.Config, log *slog.Logger) (*Service, error) {
	s := &Service{
		logger:  log,
		closers: []io.Closer{},
	}

	// Create config for the monitor
	monitorConfig := &config.Config{
		BeaconNodeURL:          cfg.BeaconNodeURL,
		TrackMissed:            cfg.TrackMissed,
		FetchIntervalSec:       cfg.FetchIntervalSec,
		EthereumRPCURL:         cfg.EthereumRPCURL,
		ValidatorOptInContract: cfg.ValidatorOptInContract,
		SlackWebhookURL:        cfg.SlackWebhookURL,
		RelayURLs:              cfg.RelayURLs,
		DashboardApiUrl:        cfg.DashboardApiUrl,
	}

	s.logger.Debug("creating duty monitor", "config", monitorConfig)
	monitor, err := monitor.New(monitorConfig, log)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	healthChecker := health.New()

	monitorDone := s.startMonitor(ctx, monitor)

	healthChecker.Register(health.CloseChannelHealthCheck("MonitorService", monitorDone))

	s.closers = append(s.closers, channelCloser(monitorDone))

	if monitor.GetDB() != nil {
		s.closers = append(s.closers, monitor.GetDB())
	}
	
	s.logger.Info("duty monitor service started",
		"beacon_node", cfg.BeaconNodeURL,
		"relay_count", len(cfg.RelayURLs))

	return s, nil
}

// startMonitor starts the monitor in a goroutine and returns a done channel
func (s *Service) startMonitor(ctx context.Context, monitor *monitor.DutyMonitor) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)
		monitor.Start(ctx)
		s.logger.Debug("monitor routine exited")
	}()

	return done
}

// Close stops the service and all its components
func (s *Service) Close() error {
	s.logger.Info("shutting down duty monitor service")

	// Cancel the context
	s.cancel()

	// Close all closers
	for _, c := range s.closers {
		if err := c.Close(); err != nil {
			return err
		}
	}

	s.logger.Info("duty monitor service shut down successfully")
	return nil
}

// channelCloser is a helper type that implements io.Closer for a channel
type channelCloser <-chan struct{}

// Close waits for the channel to be closed or times out
func (c channelCloser) Close() error {
	select {
	case <-c:
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("timed out waiting for channel to close")
	}
}
