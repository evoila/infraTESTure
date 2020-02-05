package actions

import (
	"fmt"
	"github.com/evoila/infraTESTure/config"
	"github.com/evoila/infraTESTure/infrastructure"
	"github.com/evoila/infraTESTure/infrastructure/bosh"
	"github.com/urfave/cli"
	"log"
	"plugin"
)

func Offline(c *cli.Context) error {
	conf, err := config.LoadConfig(c.String("config"))

	if err != nil {
		logError(err, "")
		return err
	}

	bosh.InitInfrastructureValues(conf)

	log.Printf("[INFO] Loading go plugin...")
	p, err := plugin.Open(conf.PreCompiledPluginPath)

	if err != nil {
		logError(err, "Could not load go plugin")
		return err
	}

	successful := 0
	failed := 0

	// Use the plugin test "Lookup" the find and execute every test found in "GetFunctionNames"
	for _, test := range conf.Testing.Tests {
		symbol, err := p.Lookup(test.Name)

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

	return nil
}
