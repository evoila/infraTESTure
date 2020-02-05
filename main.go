package main

import (
	"github.com/evoila/infraTESTure/actions"
	"github.com/urfave/cli"
	"log"
	"os"
)

var app = cli.NewApp()

func info() {
	app.Name = "infraTESTure CLI"
	app.Usage = "CLI for using the infraTESTure framework"
	app.Authors = []*cli.Author{{"Rene Schollmeyer, evoila", "rschollmeyer@evoila.de"}}
	app.Version = "0.0.1"
}

func commands() {
	app.Commands = []*cli.Command{
		{
			Name:    "info",
			Aliases: []string{"i"},
			Usage:   "Information about what tests are enabled for what services",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "repository",
					Aliases:  []string{"r"},
					Usage:    "`URL` to the Repository from which you want to get test information",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "tag",
					Aliases:  []string{"t"},
					Usage:    "Specific `TAG` to clone from",
					Required: false,
				},
			},
			Action: func(context *cli.Context) error {
				return actions.Info(context)
			},
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Run tests based on a given configuration file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "config",
					Aliases: []string{"c"},
					Usage:   "Load configuration from `FILE` for executing tests",
				},
				&cli.BoolFlag{
					Name:    "edit",
					Aliases: []string{"e"},
					Usage:   "Tells the tool if you want to edit the test code or not",
				},
				&cli.BoolFlag{
					Name:    "override",
					Aliases: []string{"o"},
					Usage:   "Overrides an already cloned repository",
				},
			},
			Action: func(context *cli.Context) error {
				return actions.Run(context)
			},
		},
		{
			Name:    "offline",
			Aliases: []string{"o"},
			Usage:   "Run tests from a pre compiled go plugin to skip git usage.",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "config",
					Aliases:  []string{"c"},
					Usage:    "Load configuration from `FILE` for executing tests",
					Required: true,
				},
			},
			Action: func(context *cli.Context) error {
				return actions.Offline(context)
			},
		},
	}
}

func main() {
	info()
	commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}
}
