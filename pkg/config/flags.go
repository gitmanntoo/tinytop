package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gitmanntoo/tinytop/pkg/core"
)

// ParseFlags parses command-line arguments and returns a Config
func ParseFlags() (core.Config, error) {
	var (
		interval    = flag.Duration("i", 1*time.Second, "collection interval (e.g., 1s, 500ms, 2m)")
		duration    = flag.Duration("d", 1*time.Second, "collection duration (e.g., 1s, 30s, 5m)")
		logFile     = flag.String("log", "tinytop.log", "log file path")
		logToStdout = flag.Bool("stdout", false, "output logs to terminal in addition to log file")
		help        = flag.Bool("h", false, "show help")
	)

	// Set custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A simple system monitoring tool.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -i 500ms -d 30s\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --interval 2s --duration 1m\n", os.Args[0])
	}

	// Add long flag support
	flag.Var((*durationValue)(interval), "interval", "collection interval (same as -i)")
	flag.Var((*durationValue)(duration), "duration", "collection duration (same as -d)")
	flag.BoolVar(help, "help", false, "show help (same as -h)")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Validate arguments
	if *interval <= 0 {
		return core.Config{}, fmt.Errorf("interval must be positive, got %v", *interval)
	}
	if *duration <= 0 {
		return core.Config{}, fmt.Errorf("duration must be positive, got %v", *duration)
	}

	return core.Config{
		Interval:    *interval,
		Duration:    *duration,
		LogFile:     *logFile,
		LogToStdout: *logToStdout,
	}, nil
}

// durationValue implements flag.Value interface for duration flags
type durationValue time.Duration

func (d *durationValue) Set(s string) error {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = durationValue(duration)
	return nil
}

func (d *durationValue) String() string {
	return time.Duration(*d).String()
}
