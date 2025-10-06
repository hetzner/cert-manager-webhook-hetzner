package main

import (
	"log/slog"
	"os"
	"strings"
)

func newLogger() *slog.Logger {
	logLevel := parseLogLevel(os.Getenv("LOG_LEVEL"))

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}))

	return logger
}

func parseLogLevel(lvl string) slog.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
