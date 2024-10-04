package main

import (
	"errors"
	"flag"
	"os"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/app"
	"github.com/Alan-Daniels/web/internal/config"
	"github.com/Alan-Daniels/web/internal/database"
)

func main() {
	staticDir := flag.String("static", ".", "where static files & default files are")
	stateDir := flag.String("state", "./tmp", "where generated files & content files are")
	configFile := flag.String("config", "./default.yml", "see config.Init()")
	flag.Parse()

	RootDir = *stateDir

	if err := InitLogger(); err != nil {
		panic(errors.Join(errors.New("Cannot start without the logger: "), err))
	}
	if config, err := config.Init(*configFile); err != nil {
		Logger.Fatal().Err(err).Str("file", *configFile).Msg("Failed to read configs")
	} else {
		Config = config
		Logger.Debug().Any("cfg", config).Msg("starting with config")
	}
	if db, err := database.Init(Config); err != nil {
		Logger.Fatal().Err(err).Msg("Failed to init the database")
	} else {
		Database = db
		Logger.Debug().Msg("connected to database")
	}

	if _, err := os.Stat(RootDir + "init.md"); errors.Is(err, os.ErrNotExist) {
		// do some first-time setup
		Logger.Warn().Str("static dir", *staticDir).Msg("Needing to do first-time setup, but it's not implimented yet!\nCrashes likely.")
	}

	if err := app.Init(); err != nil {
		Logger.Fatal().Err(err).Msg("Failed to start app")
	}
}
