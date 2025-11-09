package utils

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

const TIME_FORMAT = "2006-01-02T15:04:05.000Z07:00"

var Log zerolog.Logger

// InitLogger initializes the global logger with file and optional stdout output
// The log file is overwritten by default
func InitLogger(logFile string, logToStdout bool) error {
	// Open log file (create or truncate)
	file, err := os.Create(logFile)
	if err != nil {
		return err
	}

	// Set global time format for zerolog to RFC3339 with millisecond precision
	zerolog.TimeFieldFormat = TIME_FORMAT
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC().Truncate(time.Millisecond)
	}

	// Set up writers
	var writers []io.Writer
	writers = append(writers, file)

	if logToStdout {
		// Use console writer for human-readable terminal output
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: TIME_FORMAT,
			TimeLocation: time.UTC,
		}
		writers = append(writers, consoleWriter)
	}

	// Create multi-writer
	multi := io.MultiWriter(writers...)

	// Initialize global logger
	Log = zerolog.New(multi).With().Timestamp().Logger()

	return nil
}
