package logger

import (
	"io"
	"log"
	"log/slog"
	"os"
)

const (
	LogLevelDebug = "DEBUG"
	LogLevelInfo  = "INFO"
	LogLevelWarn  = "WARN"
	LogLevelError = "ERROR"
)

type Logger struct {
	log *slog.Logger
}

func New(logLevel string, discradMode bool) *Logger {
	var log *slog.Logger

	var w io.Writer
	level := mustParseLogLevel(logLevel)
	opts := &slog.HandlerOptions{Level: level}

	if discradMode {
		w = newDiscardWriter()
	} else {
		w = os.Stdout
	}

	log = slog.New(slog.NewTextHandler(w, opts))

	return &Logger{log: log}
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.log.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log.Warn(msg, args...)
}

func (l *Logger) Error(msg string, err error, args ...any) {
	a := []any{"error", err.Error()}
	a = append(a, args...)
	l.log.Error(msg, a...)
}

func (l *Logger) Fatal(msg string, err error, args ...any) {
	l.Error(msg, err, args...)
	os.Exit(1)
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
		log.Fatalf("could parse log level, got: %s", logLevel)
	}

	return level
}
