package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/twiny/screenshot/cmd/server"

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
			{
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
					return server.Start(c.String("config"))
				},
			},
		},
	}

	// Run CLI
	log.Println(cmd.Run(os.Args))
}
