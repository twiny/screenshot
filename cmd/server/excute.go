package server

import (
	"github.com/twiny/screenshot/api"
	"github.com/twiny/screenshot/internal/db"
	"github.com/twiny/screenshot/internal/logger"
	"github.com/twiny/screenshot/internal/rate"
)

// Start
func Start(path string) error {
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
