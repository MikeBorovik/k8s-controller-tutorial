package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected zerolog.Level
	}{
		{"trace", zerolog.TraceLevel},
		{"debug", zerolog.DebugLevel},
		{"info", zerolog.InfoLevel},
		{"warn", zerolog.WarnLevel},
		{"error", zerolog.ErrorLevel},
		{"none", zerolog.Disabled},
		{"unknown", zerolog.InfoLevel},
		{"INFO", zerolog.InfoLevel},
	}
	for _, tt := range tests {
		got := parseLogLevel(tt.input)
		if got != tt.expected {
			t.Errorf("parseLogLevel(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestInitializeLogger_JSON(t *testing.T) {
	origLogLevel := logLevel
	origLogFormat := logFormat
	defer func() {
		logLevel = origLogLevel
		logFormat = origLogFormat
	}()

	logLevel = "debug"
	logFormat = "json"

	var buf bytes.Buffer
	origStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	initializeLogger()

	// Write a debug log to trigger output
	log.Debug().Msg("test json log")

	w.Close()
	os.Stderr = origStderr
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, `"level":"debug"`) || !strings.Contains(output, `"message":"Logger initialized"`) {
		t.Errorf("Expected JSON log output, got: %s", output)
	}
}

func TestInitializeLogger_Console(t *testing.T) {
	origLogLevel := logLevel
	origLogFormat := logFormat
	defer func() {
		logLevel = origLogLevel
		logFormat = origLogFormat
	}()

	logLevel = "debug"
	logFormat = "console"

	var buf bytes.Buffer
	origStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	initializeLogger()

	log.Debug().Msg("test console log")

	w.Close()
	os.Stderr = origStderr
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "DBG") || !strings.Contains(output, "Logger initialized") {
		t.Errorf("Expected console log output, got: %s", output)
	}
}

func TestInitializeLogger_TraceLevelIncludesCaller(t *testing.T) {
	origLogLevel := logLevel
	origLogFormat := logFormat
	defer func() {
		logLevel = origLogLevel
		logFormat = origLogFormat
	}()

	logLevel = "trace"
	logFormat = "json"

	var buf bytes.Buffer
	origStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	initializeLogger()

	log.Trace().Msg("test trace log")

	w.Close()
	os.Stderr = origStderr
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, `"level":"trace"`) || !strings.Contains(output, `"caller":`) {
		t.Errorf("Expected trace log with caller, got: %s", output)
	}
}
