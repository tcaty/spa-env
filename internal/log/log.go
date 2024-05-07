package log

import (
	"fmt"
	"io"
	l "log"
	"log/slog"
	"os"
)

const (
	LogLevelDebug = "DEBUG"
	LogLevelInfo  = "INFO"
	LogLevelWarn  = "WARN"
	LogLevelError = "ERROR"
)

var log *slog.Logger

func Init(logLevel string, discradMode bool) {
	if log != nil {
		return
	}

	var w io.Writer
	level := mustParseLogLevel(logLevel)
	opts := &slog.HandlerOptions{Level: level}

	if discradMode {
		w = newDiscardWriter()
	} else {
		w = os.Stdout
	}

	log = slog.New(slog.NewTextHandler(w, opts))
}

func Debug(msg string, args ...any) {
	log.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	log.Warn(msg, args...)
}

func Error(msg string, err error, args ...any) {
	a := []any{"error", err.Error()}
	a = append(a, args...)
	log.Error(msg, a...)
}

func Fatal(msg string, err error, args ...any) {
	Error(msg, err, args...)
	os.Exit(1)
}

func ValidateLogLevel(logLevel string) error {
	switch logLevel {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		return nil
	default:
		return fmt.Errorf(
			"logLevel validation failed: available values: %s, %s, %s, %s, but got %s",
			LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError, logLevel,
		)
	}
}

func mustParseLogLevel(logLevel string) slog.Level {
	var level slog.Level

	switch logLevel {
	case LogLevelDebug:
		level = slog.LevelDebug
	case LogLevelInfo:
		level = slog.LevelInfo
	case LogLevelWarn:
		level = slog.LevelWarn
	case LogLevelError:
		level = slog.LevelError
	default:
		l.Fatalf("could parse log level, got: %s", logLevel)
	}

	return level
}
