package core

import (
	// "context"
	"encoding/json"
	"fmt"
	// "os"
	// "os/signal"
	// "syscall"
	"time"

	"github.com/gitmanntoo/tinytop/pkg/psutil"
	"github.com/gitmanntoo/tinytop/pkg/system"
	"github.com/gitmanntoo/tinytop/pkg/utils"
)

// Config holds the application configuration
type Config struct {
	Interval    time.Duration // Collection interval
	Duration    time.Duration // Collection duration
	LogFile     string        // Log file path
	LogToStdout bool          // Output logs to stdout in addition to file
}

// App represents the main application
type App struct {
	name   string
	config Config
}

// New creates a new App instance
func New(name string, config Config) *App {
	return &App{
		name:   name,
		config: config,
	}
}

// Run starts the application
func (a *App) Run() error {
	// Initialize logger
	if err := utils.InitLogger(a.config.LogFile, a.config.LogToStdout); err != nil {
		return err
	}

	// Check if running with elevated privileges
	isSudo := system.IsSudo()

	var user string
	if currentUser, err := system.GetCurrentUser(); err == nil {
		user = currentUser.Username
	} else {
		user = err.Error()
	}
	msg := "Running as regular user"

	if isSudo {
		// Show original user if available
		if sudoUser := system.GetSudoUser(); sudoUser != "" {
			user = sudoUser
		}
		msg = "Running with sudo privileges"
	}

	utils.Log.Info().
		Str("user", user).
		Bool("sudo", isSudo).
		Msg(msg)

	utils.Log.Info().
		Dur("interval", a.config.Interval).
		Dur("duration", a.config.Duration).
		Msg("Collection interval and duration")

	// Get static system information.
	sysInfo, err := psutil.Info()
	if err != nil {
		utils.Log.Error().Err(err).Msg("Failed to get system info")
		return nil
	} 

	json, err := json.MarshalIndent(sysInfo, "", "  ")
	if err != nil {
		utils.Log.Error().Err(err).Msg("Failed to marshal system info to JSON")
		return nil
	}

	fmt.Println(string(json))

	// // Set up signal handling for graceful shutdown
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// go func() {
	// 	<-sigChan
	// 	utils.Log.Info().Msg("Received interrupt signal, exiting...")
	// 	cancel()
	// }()

	// endTime := time.Now().Add(a.config.Duration)
	// currentTime := time.Now()
	// for currentTime.Before(endTime) {
	// 	// Check if context was cancelled (Ctrl-C pressed)
	// 	select {
	// 	case <-ctx.Done():
	// 		return nil
	// 	default:
	// 	}

	// 	// Get CPU times
	// 	if cpuTimes, err := cpu.Times(true); err != nil {
	// 		utils.Log.Error().Err(err).Msg("Failed to get CPU times")
	// 	} else {
	// 		for i, cpuTime := range cpuTimes {
	// 			json, err := json.Marshal(cpuTime)
	// 			if err != nil {
	// 				utils.Log.Error().Err(err).Msgf("Failed to marshal CPU %d times to JSON", i)
	// 			} else {
	// 				utils.Log.Info().
	// 					RawJSON("cpu_times", json).
	// 					Msgf("CPU %d Times (JSON)", i)
	// 			}
	// 		}
	// 	}

	// 	currentTime = currentTime.Add(a.config.Interval)
	// 	if currentTime.After(endTime) || time.Now().After(endTime) {
	// 		break
	// 	} else {
	// 		// Sleep with context awareness
	// 		timer := time.NewTimer(time.Until(currentTime))
	// 		select {
	// 		case <-ctx.Done():
	// 			timer.Stop()
	// 			return nil
	// 		case <-timer.C:
	// 		}
	// 	}
	// }

	return nil
}
