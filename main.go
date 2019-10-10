package main

import (
	"github.com/evoila/infraTESTure/config"
	"github.com/evoila/infraTESTure/infrastructure"
	"github.com/evoila/infraTESTure/infrastructure/bosh"
	"github.com/evoila/infraTESTure/parser"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"plugin"
	"strings"
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
			Flags: []cli.Flag {
				cli.StringFlag{
					Name:        "repository, r",
					Usage:       "`URL` to the Repository from which you want to get test information",
					Required:    true,
				},
			},
			Action: func(c *cli.Context) {
				url := c.String("repository")
				tmpDir := appendSlash(os.TempDir()) + "infra-tmp"

				gitClone(url, tmpDir)

				log.Printf("The following services and tests were found in %v:\n\n", url)

				dirs, err := ioutil.ReadDir(tmpDir)
				if err != nil {
					log.Fatal(err)
				}
				for _, dir := range dirs {
					if dir.IsDir() {
						goFiles, err := ioutil.ReadDir(appendSlash(tmpDir)+dir.Name())
						if err != nil {
							log.Fatal(err)
						}

						for _, goFile := range goFiles {
							if goFile.Name() == dir.Name()+".go" {
								log.Printf("├── %v", color.GreenString(dir.Name()))
								parser.GetAnnotations(appendSlash(tmpDir)+dir.Name()+"/"+goFile.Name())
							}
						}
					}

				}

				cmd := exec.Command("bash", "-c", "rm -rf " + tmpDir)
				err = cmd.Run()

				if err != nil {
					logError(err, "Could not delete directory")
				}
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

				var repoPath string

				if repoPath = appendSlash(os.TempDir()); conf.Github.SavingLocation != "" {
					repoPath = appendSlash(conf.Github.SavingLocation)
				}

				repoPath += conf.Github.RepoName

				if _, err := os.Stat(repoPath); os.IsNotExist(err) {
					log.Printf("[INFO] Cloning repository from %v\n", conf.Github.TestRepo)
					gitClone(conf.Github.TestRepo, repoPath)
				} else {
					log.Printf("[INFO] Using existing repository %v\n", repoPath)
				}

				serviceDir := repoPath + "/" + conf.Service.Name

				if c.Bool("edit") {
					cmd := exec.Command("bash", "-c", "open -t " + appendSlash(serviceDir) + conf.Service.Name + ".go")
					err = cmd.Run()

					if err != nil {
						logError(err, "Could not open test file")
					}
				} else {
					log.Printf("[INFO] Building go plugin from directory %v\n", serviceDir)
					cmd := exec.Command("bash", "-c", "cd " + serviceDir + " && go build -buildmode=plugin")
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
				}
			},
		},
	}
}

func logError(err error, customMessage string) {
	log.Fatal(color.RedString("[ERROR] " + customMessage + ": " + err.Error()))
}

func gitClone(url string, repoPath string) {
	cmd := exec.Command("bash", "-c", "git clone " + url + " " + repoPath)
	err := cmd.Run()

	if err != nil {
		logError(err, "Could not clone repository")
	}
}

func appendSlash(dir string) string {
	if !strings.HasSuffix(dir, "/") {
		return dir + "/"
	}
	return dir
}


func main() {
	info()
	commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}
}