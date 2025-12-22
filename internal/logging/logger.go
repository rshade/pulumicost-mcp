package logging

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
)

// Logger wraps slog for structured JSON logging
type Logger struct {
	*slog.Logger
}

// Config holds logger configuration
type Config struct {
	Level  string // debug, info, warn, error
	Format string // json, text
	Output io.Writer
}

// New creates a new structured logger
func New(cfg Config) *Logger {
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	level := parseLevel(cfg.Level)

	var handler slog.Handler
	if cfg.Format == "text" {
		handler = slog.NewTextHandler(cfg.Output, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		// Default to JSON format
		handler = slog.NewJSONHandler(cfg.Output, &slog.HandlerOptions{
			Level: level,
		})
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// Default creates a logger with sensible defaults
func Default() *Logger {
	return New(Config{
		Level:  "info",
		Format: "json",
		Output: os.Stdout,
	})
}

// WithContext adds context values to the logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Extract common context values if they exist
	return l
}

// WithService adds service name to all log entries
func (l *Logger) WithService(service string) *Logger {
	if l == nil || l.Logger == nil {
		// Return a no-op logger if nil
		return &Logger{Logger: slog.New(slog.NewJSONHandler(io.Discard, nil))}
	}
	return &Logger{
		Logger: l.Logger.With("service", service),
	}
}

// WithRequest adds request ID and metadata
func (l *Logger) WithRequest(requestID string) *Logger {
	return &Logger{
		Logger: l.Logger.With("request_id", requestID),
	}
}

// InfoJSON logs structured data as JSON
func (l *Logger) InfoJSON(msg string, data interface{}) {
	bytes, _ := json.Marshal(data)
	l.Info(msg, "data", string(bytes))
}

// ErrorJSON logs error with structured data
func (l *Logger) ErrorJSON(msg string, err error, data interface{}) {
	bytes, _ := json.Marshal(data)
	l.Error(msg, "error", err, "data", string(bytes))
}

func parseLevel(level string) slog.Level {
	switch level {
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
