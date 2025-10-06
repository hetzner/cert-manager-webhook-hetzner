package main

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLogLevel(t *testing.T) {
	for value, level := range map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,

		"ERROR":   slog.LevelError,
		"unknown": slog.LevelInfo,
	} {
		assert.Equal(t, level, parseLogLevel(value))
	}
}
