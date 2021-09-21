package main

import (
	"github.com/twiny/screenshot/cmd/screen/api"
	"github.com/twiny/screenshot/internal/db"
	"github.com/twiny/screenshot/pkg/logger"
	"github.com/twiny/screenshot/pkg/rate"
)

// Start
func startServer(path string) error {
	// config
	config, err := api.YAMLConfig(path)
	if err != nil {
		return err
	}

	// store
	store, err := db.NewStore(config.StorePath)
	if err != nil {
		return err
	}

	// limiter
	limiter := rate.NewLimiter(config.Rate, config.Bursts)

	// logger
	logger, err := logger.NewLogger(config.LogPath)
	if err != nil {
		return err
	}

	// new app
	app := api.NewApp(config, store, limiter, logger)
	defer app.Close()

	// clean store & limiter
	app.Sync()

	// start http server
	return api.NewHTTPServer(app).ListenAndServe(":" + config.Port)
}