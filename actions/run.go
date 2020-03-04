package actions

import (
	"fmt"
	"github.com/evoila/infraTESTure/config"
	"github.com/evoila/infraTESTure/infrastructure"
	"github.com/evoila/infraTESTure/infrastructure/bosh"
	"github.com/evoila/infraTESTure/logger"
	"github.com/evoila/infraTESTure/parser"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"plugin"
	"runtime"
	"strings"
)

func Run(c *cli.Context) error {
	conf, err := config.LoadConfig(c.String("config"))

	if err != nil {
		logger.LogError(err, "")
		return err
	}

	bosh.InitInfrastructureValues(conf)

	var repoPath string

	if repoPath = appendSlash(os.TempDir()); conf.Github.SavingLocation != "" {
		repoPath = appendSlash(conf.Github.SavingLocation)
	}

	repoPath += conf.Github.RepoName

	if c.Bool("override") {
		cmd := exec.Command("bash", "-c", "rm -rf "+repoPath)
		err = cmd.Run()

		if err != nil {
			logger.LogError(err, "")
			return err
		}
	}

	// Check if the repository is already cloned, and if so use this repository
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		logger.LogInfoF("[INFO] Cloning repository from %v\n", conf.Github.TestRepo)
		err = gitClone(conf.Github.TestRepo, repoPath, conf.Github.Tag)
		if err != nil {
			logger.LogError(err, "Could not clone repository")
			return err
		}
	} else {
		logger.LogInfoF("[INFO] Using existing repository %v\n", repoPath)
	}

	serviceDir := repoPath + "/" + conf.Service.Name

	// If the --edit flag is set open the test file instead of running the tests
	if c.Bool("edit") {
		//TODO: change since this fails with changes in infra-tests
		cmd := exec.Command("bash", "-c", "code . "+serviceDir)
		err = cmd.Run()

		if err != nil {
			logger.LogError(err, "Could not open test file")
			return err
		}
	} else {
		// Build the given test repository as a go plugin
		logger.LogInfoF("[INFO] Building go plugin from directory %v\n", serviceDir)
		cmd := exec.Command("bash", "-c", "cd "+serviceDir+" && "+runtime.Version()+" build -buildmode=plugin")
		err = cmd.Run()

		if err != nil {
			logger.LogError(err, "Could not build go plugin")
			return err
		}

		logger.LogInfoF("[INFO] Loading go plugin...")
		p, err := plugin.Open(serviceDir + "/" + conf.Service.Name + ".so")

		if err != nil {
			logger.LogError(err, "Could not load go plugin")
			return err
		}

		var functionNames []string

		files, err := ioutil.ReadDir(serviceDir)

		if err != nil {
			logger.LogError(err, "Could not load service directory")
			return err
		}

		// Get a list of all names of functions that has to be executed, based on if their annotations match the
		// test names provided by the configuration.yml
		for _, test := range conf.Testing.Tests {
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".go") {
					newFunctionNames, err := parser.GetFunctionNames(test, appendSlash(serviceDir)+file.Name())

					if err != nil {
						logger.LogError(err, "")
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
				logger.LogError(err, "")
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

		logger.LogInfoF("[INFO] %d of %d tests succeeded", successful, successful+failed)
	}

	return nil
}
