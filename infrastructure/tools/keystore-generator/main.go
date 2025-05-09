package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var (
	optionKeystoreDir = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "keystore-dir",
		Value:   "./",
		Usage:   "Directory to save the keystore file",
		EnvVars: []string{"KEYSTOREGEN_KEYSTORE_DIR"},
	})

	optionPassphrase = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "passphrase",
		Usage:   "Passphrase for encrypting the keystore",
		EnvVars: []string{"KEYSTOREGEN_PASSPHRASE"},
	})

	optionKeystoreFile = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "keystore-file",
		Usage:   "Path to the Ethereum keystore file",
		EnvVars: []string{"KEYSTOREGEN_KEYSTORE_FILE"},
	})

	optionPrivateKey = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "private-key",
		Usage:   "Private key to import (as hex string)",
		EnvVars: []string{"KEYSTOREGEN_PRIVATE_KEY"},
	})

	optionLogLevel = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-level",
		Usage:   "Log level, options are 'debug', 'info', 'warn', 'error'",
		EnvVars: []string{"KEYSTOREGEN_LOG_LEVEL"},
		Value:   "info",
		Action: func(_ *cli.Context, s string) error {
			if !slices.Contains([]string{"debug", "info", "warn", "error"}, s) {
				return fmt.Errorf("invalid value: -log-level=%q", s)
			}
			return nil
		},
	})

	optionLogFmt = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "Log format, options are 'text' or 'json'",
		EnvVars: []string{"KEYSTOREGEN_LOG_FMT"},
		Value:   "text",
		Action: func(_ *cli.Context, s string) error {
			if !slices.Contains([]string{"text", "json"}, s) {
				return fmt.Errorf("invalid value: -log-fmt=%q", s)
			}
			return nil
		},
	})

	optionLogTags = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-tags",
		Usage:   "Log tags is a comma-separated list of <name:value> pairs that will be inserted into each log line",
		EnvVars: []string{"KEYSTOREGEN_LOG_TAGS"},
		Action: func(ctx *cli.Context, s string) error {
			for i, p := range strings.Split(s, ",") {
				if len(strings.Split(p, ":")) != 2 {
					return fmt.Errorf("invalid value at index %d: -log-tags=%q", i, s)
				}
			}
			return nil
		},
	})
)

func main() {
	flags := []cli.Flag{
		optionKeystoreDir,
		optionPassphrase,
		optionKeystoreFile,
		optionPrivateKey,
		optionLogFmt,
		optionLogLevel,
		optionLogTags,
	}
	app := &cli.App{
		Name:  "keystoreGenerator",
		Usage: "Generates an Ethereum keystore file.",
		Commands: []*cli.Command{
			{
				Name:   "generate",
				Usage:  "Generates a new Ethereum keystore file.",
				Flags:  flags,
				Action: newAction((*command).generateKeystore),
			},
			{
				Name:   "extract",
				Usage:  "Extracts the private key from a keystore file.",
				Flags:  flags,
				Action: newAction((*command).extractPrivateKey),
			},
			{
				Name:   "import",
				Usage:  "Imports a private key and creates a keystore file.",
				Flags:  flags,
				Action: newAction((*command).importPrivateKey),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		//nolint:errcheck
		fmt.Fprintln(app.ErrWriter, err)
	}
}

// newAction returns new cli.ActionFunc with injected logger.
func newAction(action func(*command, *cli.Context) error) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		logger, err := newLogger(
			ctx.String(optionLogLevel.Name),
			ctx.String(optionLogFmt.Name),
			ctx.String(optionLogTags.Name),
			ctx.App.Writer,
		)
		if err != nil {
			return fmt.Errorf("failed to create logger: %w", err)
		}
		if err := action(&command{logger}, ctx); err != nil {
			logger.Error("command execution failed", "error", err)
			return err
		}
		return nil
	}
}

// command groups application commands.
type command struct {
	logger *slog.Logger
}

func (c *command) generateKeystore(ctx *cli.Context) error {
	passphrase := ctx.String(optionPassphrase.Name)
	if passphrase == "" {
		return fmt.Errorf("-%s is required", optionPassphrase.Name)
	}

	keystoreDir := ctx.String(optionKeystoreDir.Name)
	ks := keystore.NewKeyStore(keystoreDir, keystore.LightScryptN, keystore.LightScryptP)
	account, err := ks.NewAccount(passphrase)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	c.logger.Info("new keystore generated successfully", "address", account.Address.Hex(), "path", account.URL.Path)
	return nil
}

func (c *command) extractPrivateKey(ctx *cli.Context) error {
	keystoreFile := ctx.String(optionKeystoreFile.Name)
	if keystoreFile == "" {
		return fmt.Errorf("-%s is required", optionKeystoreFile.Name)
	}

	passphrase := ctx.String(optionPassphrase.Name)
	if passphrase == "" {
		return fmt.Errorf("-%s is required", optionPassphrase.Name)
	}

	keyjson, err := os.ReadFile(keystoreFile)
	if err != nil {
		return fmt.Errorf("failed to read keystore file: %w", err)
	}

	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		return fmt.Errorf("failed to decrypt key: %w", err)
	}

	c.logger.Info("private key extracted successfully", "key", fmt.Sprintf("%x", key.PrivateKey.D.Bytes()))
	return nil
}

func (c *command) importPrivateKey(ctx *cli.Context) error {
	privateKeyHex := ctx.String(optionPrivateKey.Name)
	if privateKeyHex == "" {
		return fmt.Errorf("-%s is required", optionPrivateKey.Name)
	}

	passphrase := ctx.String(optionPassphrase.Name)
	if passphrase == "" {
		return fmt.Errorf("-%s is required", optionPassphrase.Name)
	}

	// Convert hex string to ECDSA private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to decode private key: %v", err)
	}

	// Create keystore and import the private key
	keystoreDir := ctx.String(optionKeystoreDir.Name)
	ks := keystore.NewKeyStore(keystoreDir, keystore.LightScryptN, keystore.LightScryptP)
	account, err := ks.ImportECDSA(privateKey, passphrase)
	if err != nil {
		return fmt.Errorf("failed to import private key: %v", err)
	}

	c.logger.Info("account imported successfully", "address", account.Address.Hex(), "path", account.URL.Path)
	return nil
}

// newLogger returns a new instance of slog.Logger configured with the given parameters.
func newLogger(lvl, logFmt, tags string, sink io.Writer) (*slog.Logger, error) {
	level := new(slog.LevelVar)
	if err := level.UnmarshalText([]byte(lvl)); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	var (
		handler slog.Handler
		options = &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		}
	)
	switch logFmt {
	case "text":
		handler = slog.NewTextHandler(sink, options)
	case "json", "none":
		handler = slog.NewJSONHandler(sink, options)
	default:
		return nil, fmt.Errorf("invalid log format: %s", logFmt)
	}

	logger := slog.New(handler)

	if tags == "" {
		return logger, nil
	}

	var args []any
	for i, p := range strings.Split(tags, ",") {
		kv := strings.Split(p, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid tag at index %d", i)
		}
		args = append(args, strings.ToValidUTF8(kv[0], "�"), strings.ToValidUTF8(kv[1], "�"))
	}
	return logger.With(args...), nil
}
