package bosh

import (
	"github.com/evoila/infraTESTure/config"
	"github.com/evoila/infraTESTure/logger"
	"github.com/fatih/color"
	"strconv"
)

type Bosh struct{}

var conf *config.Config

// Initialize the Bosh Director and the Deployment affiliated to the deployment name in the config
// @param config Configuration struct from github.com/evoila/infraTESTure/config
func InitInfrastructureValues(config *config.Config) {
	conf = config

	deploymentName = config.DeploymentName

	var err error

	// Create a bosh director to get the deployment
	boshDirector, err = BuildDirector(config)
	if err != nil {
		logError(err, "")
	}

	// See director/interfaces.go for a full list of methods.
	deployment, err = boshDirector.FindDeployment(config.DeploymentName)
	if err != nil {
		logError(err, "")
	}
}

func logError(err error, customMessage string) {
	if err != nil {
		logger.LogErrorF(color.RedString("[ERROR] " + customMessage + ": " + err.Error()))
	} else {
		logger.LogErrorF(color.RedString("[ERROR] " + customMessage))
	}
}

func stringToFloat(value string) float64 {
	if value == "" {
		return 0
	}

	floatValue, err := strconv.ParseFloat(value, 64)

	if err != nil {
		logError(err, "")
	}

	return floatValue
}
