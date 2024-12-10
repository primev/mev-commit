package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/primev/mev-commit/cl/streamer"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var (
	configFlag = &cli.StringFlag{
		Name:    "config",
		Usage:   "Path to config file",
		EnvVars: []string{"STREAMER_CONFIG"},
	}

	redisAddrFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "redis-addr",
		Usage:   "Redis address",
		EnvVars: []string{"STREAMER_REDIS_ADDR"},
		Value:   "127.0.0.1:7001",
		Action: func(_ *cli.Context, s string) error {
			host, port, err := net.SplitHostPort(s)
			if err != nil {
				return fmt.Errorf("invalid redis-addr: %v", err)
			}
			if host == "" {
				return fmt.Errorf("invalid redis-addr: missing host")
			}
			if p, err := strconv.Atoi(port); err != nil || p <= 0 || p > 65535 {
				return fmt.Errorf("invalid redis-addr: invalid port number")
			}
			return nil
		},
	})

	listenAddrFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "listen-addr",
		Usage:   "Streamer listen address",
		EnvVars: []string{"STREAMER_LISTEN_ADDR"},
		Value:   ":50051",
	})

	logFmtFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "Log format to use, options are 'text' or 'json'",
		EnvVars: []string{"STREAMER_LOG_FMT"},
		Value:   "text",
	})

	logLevelFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-level",
		Usage:   "Log level to use, options are 'debug', 'info', 'warn', 'error'",
		EnvVars: []string{"STREAMER_LOG_LEVEL"},
		Value:   "info",
	})
)

type Config struct {
	RedisAddr  string
	ListenAddr string
}

func main() {
	flags := []cli.Flag{
		configFlag,
		redisAddrFlag,
		listenAddrFlag,
		logFmtFlag,
		logLevelFlag,
	}

	app := &cli.App{
		Name:  "streamer",
		Usage: "Start the streamer",
		Flags: flags,
		Before: altsrc.InitInputSourceWithContext(flags,
			func(c *cli.Context) (altsrc.InputSourceContext, error) {
				configFile := c.String("config")
				if configFile != "" {
					return altsrc.NewYamlSourceFromFile(configFile)
				}
				return &altsrc.MapInputSource{}, nil
			}),
		Action: func(c *cli.Context) error {
			return startStreamer(c)
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error running streamer:", err)
	}
}

func startStreamer(c *cli.Context) error {
	log, err := util.NewLogger(
		c.String(logLevelFlag.Name),
		c.String(logFmtFlag.Name),
		"", // No log tags
		c.App.Writer,
	)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Load configuration
	cfg := Config{
		RedisAddr:  c.String(redisAddrFlag.Name),
		ListenAddr: c.String(listenAddrFlag.Name),
	}

	log.Info("Starting streamer with configuration", "config", cfg)

	// Initialize the Streamer
	streamer, err := streamer.NewPayloadStreamer(cfg.RedisAddr, log)
	if err != nil {
		log.Error("Failed to initialize Streamer", "error", err)
		return err
	}

	ctx, stop := signal.NotifyContext(c.Context, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start the streamer
	go func() {
		if err := streamer.Start(cfg.ListenAddr); err != nil {
			log.Error("Streamer exited with error", "error", err)
			stop()
		}
	}()

	<-ctx.Done()

	streamer.Stop()

	log.Info("Streamer shutdown completed")
	return nil
}
