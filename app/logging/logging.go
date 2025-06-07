package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config holds logging configuration
type Config struct {
	Level      string
	Format     string
	Output     string
	TimeFormat string
}

// Setup initializes the logger
func Setup(cfg Config) error {
	// Set log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	zerolog.SetGlobalLevel(level)

	// Set time format
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = time.RFC3339
	}
	zerolog.TimeFieldFormat = cfg.TimeFormat

	// Set output format
	var output zerolog.ConsoleWriter
	if cfg.Format == "json" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: cfg.TimeFormat,
		}
	} else {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: cfg.TimeFormat,
			NoColor:    false,
		}
	}

	// Set global logger
	log.Logger = zerolog.New(output).With().Timestamp().Logger()

	return nil
}

// Logger wraps zerolog.Logger with additional methods
type Logger struct {
	log zerolog.Logger
}

// New creates a new logger
func New() *Logger {
	return &Logger{
		log: log.Logger,
	}
}

// With creates a new logger with the given fields
func (l *Logger) With(fields map[string]interface{}) *Logger {
	return &Logger{
		log: l.log.With().Fields(fields).Logger(),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	l.log.Debug().Fields(fields).Msg(msg)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	l.log.Info().Fields(fields).Msg(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	l.log.Warn().Fields(fields).Msg(msg)
}

// Error logs an error message
func (l *Logger) Error(err error, msg string, fields map[string]interface{}) {
	l.log.Error().Err(err).Fields(fields).Msg(msg)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(err error, msg string, fields map[string]interface{}) {
	l.log.Fatal().Err(err).Fields(fields).Msg(msg)
}

// WithError adds an error to the logger
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		log: l.log.With().Err(err).Logger(),
	}
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		log: l.log.With().Interface(key, value).Logger(),
	}
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	return &Logger{
		log: l.log.With().Fields(fields).Logger(),
	}
}
