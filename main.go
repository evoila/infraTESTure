package main

import (
	"github.com/evoila/infraTESTure/config"
	"github.com/evoila/infraTESTure/infrastructure"
	"github.com/evoila/infraTESTure/infrastructure/bosh"
	"github.com/evoila/infraTESTure/parser"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/exec"
	"plugin"
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
				//TODO: infra-tests directory tree print
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
				conf, err := config.LoadConfig(c.String("config"))

				if err != nil {
					logError(err, "")
				}

				bosh.InitInfrastructureValues(conf)

				//TODO: Check if repository already exists
				log.Printf("[INFO] Cloning repository from %v\n", conf.TestRepo)
				cmd := exec.Command("bash", "-c", "git clone " + conf.TestRepo + " .tmp/infra-tests")
				err = cmd.Run()

				if err != nil {
					logError(err, "Could not clone repository")
				}

				serviceDir := ".tmp/infra-tests/" + conf.Service.Name

				log.Printf("[INFO] Building go plugin from directory %v\n", serviceDir)
				cmd = exec.Command("bash", "-c", "cd " + serviceDir + " && go build -buildmode=plugin")
				err = cmd.Run()

				if err != nil {
					logError(err, "Could not build go plugin")
				}

				log.Printf("[INFO] Loading go plugin...")
				p, err := plugin.Open(serviceDir + "/" + conf.Service.Name + ".so")

				if err != nil {
					logError(err, "Could not load go plugin")
				}

				methodNames, err := parser.GetMethodNames(conf.Testing.Tests, serviceDir + "/" + conf.Service.Name + ".go")

				if err != nil {
					logError(err, "")
				}

				for _, method := range methodNames {
					symbol, err := p.Lookup(method)

					if err != nil {
						logError(err, "")
					}

					fun, ok := symbol.(func(*config.Config, infrastructure.Infrastructure))

					if !ok {
						panic(ok)
					}

					fun(conf, bosh.Bosh{})
				}
			},
		},
	}
}

func logError(err error, customMessage string) {
	log.Fatal(color.RedString("[ERROR] " + customMessage + ": " + err.Error()))
}

func main() {
	info()
	commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}
}