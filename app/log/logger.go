package log

import (
	"log/slog"
	"os"
	"sync"

	"warhoop/app/config"
)

type Logger struct {
	*slog.Logger
	Uptrace *UptraceLogger
}

var (
	logger *Logger
	once   sync.Once
)

// Get initializes logger with log level from config. Once.
func Get() *Logger {
	once.Do(func() {
		// Get config
		cfg := config.Get()
		level := slog.LevelInfo

		if cfg != nil {
			switch cfg.Service.LogLevel {
			case "debug":
				level = slog.LevelDebug
			case "info":
				level = slog.LevelInfo
			case "warn", "warning":
				level = slog.LevelWarn
			case "err", "error":
				level = slog.LevelError
			case "fatal":
				level = slog.LevelError
			case "panic":
				level = slog.LevelError
			default:
				level = slog.LevelInfo
			}
		}

		consoleOpts := &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		}
		consoleHandler := slog.NewTextHandler(os.Stdout, consoleOpts)

		uptraceLogger := NewUptraceLogger()

		// Create logger
		logger = &Logger{
			Logger:  slog.New(consoleHandler),
			Uptrace: uptraceLogger,
		}
	})

	return logger
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.Logger.Debug(msg, fieldsToAny(fields)...)
	l.Uptrace.Debug(msg, fields)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.Logger.Info(msg, fieldsToAny(fields)...)
	l.Uptrace.Info(msg, fields)
}

func fieldsToAny(fields []Field) []any {
	result := make([]any, len(fields))
	for i, field := range fields {
		result[i] = field
	}
	return result
}
