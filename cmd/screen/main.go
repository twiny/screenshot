package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/twiny/screenshot/cmd/screen/api"
	"github.com/twiny/screenshot/internal/db"
	"github.com/twiny/screenshot/pkg/logger"
	"github.com/twiny/screenshot/pkg/rate"
	"github.com/urfave/cli/v2"
)

//go:embed version
var version string

// main
func main() {
	cmd := &cli.App{
		Name:     "screenshot",
		HelpName: "screenshot",
		Usage:    "Capture webpage screenshot",
		Version:  version,
		Commands: []*cli.Command{
			startServerCmd(),
		},
	}

	// Run CLI
	log.Println(cmd.Run(os.Args))
}

// startServerCmd
func startServerCmd() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Start HTTP server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "path to config.yaml `PATH`",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			return startServer(c.String("config"))
		},
	}
}

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
