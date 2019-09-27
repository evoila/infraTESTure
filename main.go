package main

import (
	"fmt"
	"github.com/evoila/infraTESTure/bosh"
	"github.com/evoila/infraTESTure/config"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"log"
	"os"
)

var app = cli.NewApp()

func info() {
	app.Name = "infraTESTure CLI"
	app.Email = "rschollmeyer@evoila.de"
	app.Usage = "CLI for using the infraTESTure framework"
	app.Author = "Rene Schollmeyer, evoila"
	app.Version = "0.0.1"
}

func commands() {
	app.Commands = []cli.Command {
		{
			Name: "info",
			Aliases: []string{"i"},
			Usage: "Information about what tests are enabled for what services",
			Action: func(c *cli.Context) {
				//TODO: Implementieren
			},
		},
		{
			Name: "run",
			Aliases: []string{"r"},
			Usage: "Run tests based on a given configuration file",
			Flags: []cli.Flag {
				cli.StringFlag{
					Name:        "config, c",
					Usage:       "Load configuration from `FILE` for executing tests",
				},
				cli.BoolFlag{
					Name:        "edit, e",
					Usage:       "Tells the tool if you want to edit the test code or not",
				},
			},
			Action: func(c *cli.Context) {
				config, err := config.LoadConfig(c.String("config"))

				if err != nil {
					log.Fatal(err)
				}

				for _, test := range config.Testing.Tests {
					switch test.Name {
					case "health":
						health := bosh.IsDeploymentRunning(config)
						fmt.Printf("\nDeployment %v is ", config.DeploymentName)
						if  health {
							color.Green("healthy")
						} else {
							color.Red("not healthy")
						}
					case "service":
						log.Println("Coming soon...")
						//TODO: Implementieren
					case "failover":
						log.Println("Coming soon...")
						//TODO: Implementieren
					}

					fmt.Printf("\n##########\n\n")
				}
			},
		},
	}
}

func main() {
	info()
	commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}