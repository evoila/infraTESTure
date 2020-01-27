package main

import (
	"fmt"
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
	"sort"
	"strings"
)

var app = cli.NewApp()

func info() {
	app.Name = "infraTESTure CLI"
	app.Usage = "CLI for using the infraTESTure framework"
	app.Authors = []*cli.Author{{"Rene Schollmeyer, evoila", "rschollmeyer@evoila.de"}}
	app.Version = "0.0.1"
}

func commands() {
	app.Commands = []*cli.Command {
		{
			Name: "info",
			Aliases: []string{"i"},
			Usage: "Information about what tests are enabled for what services",
			Flags: []cli.Flag {
				&cli.StringFlag{
					Name:        "repository",
					Aliases:	 []string{"r"},
					Usage:       "`URL` to the Repository from which you want to get test information",
					Required:    true,
				},
				&cli.StringFlag{
					Name:        "tag",
					Aliases:	 []string{"t"},
					Usage:       "Specific `TAG` to clone from",
					Required:    false,
				},
			},
			Action: func(c *cli.Context) error {
				url := c.String("repository")
				tag := c.String("tag")
				tmpDir := appendSlash(os.TempDir()) + "infra-tmp"

				err := gitClone(url, tmpDir, tag)

				if err != nil {
					logError(err, "Could not clone repository")
					return err
				}

				log.Printf("The following services and tests were found in %v:\n\n", url)

				dirs, err := ioutil.ReadDir(tmpDir)
				if err != nil {
					logError(err, "No repository found")
					return err
				}

				// Iterate through all directories and files inside these directories to get a list
				// of all annotations, and therefore a list of all offered tests
				for _, dir := range dirs {
					if dir.IsDir() {
						goFiles, err := ioutil.ReadDir(appendSlash(tmpDir)+dir.Name())
						if err != nil {
							logError(err, "Failed to acquire offered tests")
							return err
						}

						if !strings.HasPrefix(dir.Name(), ".") {
							log.Printf("├── %v", color.GreenString(dir.Name()))
						}

						var testNames []string

						for _, goFile := range goFiles {
							if strings.HasSuffix(goFile.Name(), ".go") {
								tmpTestNames := parser.GetAnnotations(appendSlash(tmpDir)+dir.Name()+"/"+goFile.Name())

								for i := range tmpTestNames {
									testNames = append(testNames, tmpTestNames[i])
								}
							}
						}

						sort.Strings(testNames)

						for i := range testNames {
							log.Printf("│ \t├── %v", testNames[i])
						}
					}
				}

				cmd := exec.Command("bash", "-c", "rm -rf " + tmpDir)
				err = cmd.Run()

				if err != nil {
					logError(err, "Could not delete directory")
					return err
				}

				return nil
			},
		},
		{
			Name: "run",
			Aliases: []string{"r"},
			Usage: "Run tests based on a given configuration file",
			Flags: []cli.Flag {
				&cli.StringFlag{
					Name:        "config",
					Aliases:	 []string{"c"},
					Usage:       "Load configuration from `FILE` for executing tests",
				},
				&cli.BoolFlag{
					Name:        "edit",
					Aliases:	 []string{"e"},
					Usage:       "Tells the tool if you want to edit the test code or not",
				},
				&cli.BoolFlag{
					Name: 		 "override",
					Aliases:	 []string{"o"},
					Usage:		 "Overrides an already cloned repository",
				},
			},

			Action: func(c *cli.Context) error {
				conf, err := config.LoadConfig(c.String("config"))

				if err != nil {
					logError(err, "")
					return err
				}

				bosh.InitInfrastructureValues(conf)

				var repoPath string

				if repoPath = appendSlash(os.TempDir()); conf.Github.SavingLocation != "" {
					repoPath = appendSlash(conf.Github.SavingLocation)
				}

				repoPath += conf.Github.RepoName

				if c.Bool("override") {
					cmd := exec.Command("bash", "-c", "rm -rf " + repoPath)
					err = cmd.Run()

					if err != nil {
						logError(err, "")
						return err
					}
				}

				// Check if the repository is already cloned, and if so use this repository
				if _, err := os.Stat(repoPath); os.IsNotExist(err) {
					log.Printf("[INFO] Cloning repository from %v\n", conf.Github.TestRepo)
					err = gitClone(conf.Github.TestRepo, repoPath, conf.Github.Tag)
					if err != nil {
						logError(err, "Could not clone repository")
						return err
					}
				} else {
					log.Printf("[INFO] Using existing repository %v\n", repoPath)
				}

				serviceDir := repoPath + "/" + conf.Service.Name

				// If the --edit flag is set open the test file instead of running the tests
				if c.Bool("edit") {
					//TODO: change since this fails with changes in infra-tests
					cmd := exec.Command("bash", "-c", "code . " + serviceDir)
					err = cmd.Run()

					if err != nil {
						logError(err, "Could not open test file")
						return err
					}
				} else {
					// Build the given test repository as a go plugin
					log.Printf("[INFO] Building go plugin from directory %v\n", serviceDir)
					cmd := exec.Command("bash", "-c", "cd " + serviceDir + " && go build -buildmode=plugin")
					err = cmd.Run()

					if err != nil {
						logError(err, "Could not build go plugin")
						return err
					}

					log.Printf("[INFO] Loading go plugin...")
					p, err := plugin.Open(serviceDir + "/" + conf.Service.Name + ".so")

					if err != nil {
						logError(err, "Could not load go plugin")
						return err
					}

					var functionNames []string

					files, err := ioutil.ReadDir(serviceDir)

					if err != nil {
						logError(err, "Could not load service directory")
						return err
					}

					// Get a list of all names of functions that has to be executed, based on if their annotations match the
					// test names provided by the configuration.yml
					for _, test := range conf.Testing.Tests {
						for _, file := range files {
							if strings.HasSuffix(file.Name(), ".go") {
								newFunctionNames, err := parser.GetFunctionNames(test, appendSlash(serviceDir) + file.Name())

								if err != nil {
									logError(err, "")
									return err
								}

								functionNames = append(functionNames, newFunctionNames...)
							}
						}
					}

					successful := 0
					failed := 0

					// Use the plugin function "Lookup" the find and execute every function found in "GetFunctionNames"
					for _, function := range functionNames {
						symbol, err := p.Lookup(function)

						if err != nil {
							logError(err, "")
							return err
						}

						fun, ok := symbol.(func(*config.Config, infrastructure.Infrastructure) bool)

						if !ok {
							panic(ok)
						}

						testResult := fun(conf, &bosh.Bosh{})

						if testResult {
							successful++
						} else {
							failed++
						}
					}

					fmt.Printf("\033[1;34m%s\033[0m", "\n##### Result #####\n")

					log.Printf("[INFO] %d of %d tests succeeded", successful, successful+failed)
				}

				return nil
			},
		},
	}
}

func logError(err error, customMessage string) {
	log.Printf(color.RedString("[ERROR] " + customMessage + ": " + err.Error()))
}

func gitClone(url string, repoPath string, tag string) error {
	var tagClone string

	if tag != "" {
		tagClone = "--branch " + tag + " --single-branch"
	}

	cmd := exec.Command("bash", "-c", "git clone " + url + " " + repoPath + " " + tagClone)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
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