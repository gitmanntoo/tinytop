package main

import (
	"os"

	"github.com/gitmanntoo/tinytop/pkg/config"
	"github.com/gitmanntoo/tinytop/pkg/core"
	"github.com/gitmanntoo/tinytop/pkg/utils"
)

func main() {
	// Parse command-line flags
	cfg, err := config.ParseFlags()
	if err != nil {
		// Use stderr for early errors before logger is initialized
		utils.Log.Fatal().Err(err).Msg("Configuration error")
	}

	app := core.New("tinytop", cfg)

	if err := app.Run(); err != nil {
		utils.Log.Fatal().Err(err).Msg("Application error")
		os.Exit(1)
	}
}
