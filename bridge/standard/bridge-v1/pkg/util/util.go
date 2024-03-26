package util

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
)

// NewLogger initializes a *slog.Logger with specified level, format, and sink.
//   - lvl: string representation of slog.Level
//   - logFmt: format of the log output: "text", "json", "none" defaults to "json"
//   - tags: comma-separated list of <name:value> pairs that will be inserted into each log line
//   - sink: destination for log output (e.g., os.Stdout, file)
//
// Returns a configured *slog.Logger on success or nil on failure.
func NewLogger(lvl, logFmt, tags string, sink io.Writer) (*slog.Logger, error) {
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
