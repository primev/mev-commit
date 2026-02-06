package logutil

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func SetLogLevel(lvl string, format string, tags map[string]string) error {
	level := new(slog.LevelVar)
	if err := level.UnmarshalText([]byte(lvl)); err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	var handler slog.Handler

	switch format {
	case "json":
		// JSON format
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		})
	case "text":
		// Text format (default)
		opts := &tint.Options{
			Level:      level,
			TimeFormat: time.Kitchen, // Optional: Customize the time format
			AddSource:  true,
		}
		handler = tint.NewHandler(os.Stdout, opts)
	}
	var logTags []any
	for k, v := range tags {
		logTags = append(logTags, k, v)
	}
	logger := slog.New(handler).With(logTags...)
	slog.SetDefault(logger)
	return nil
}
