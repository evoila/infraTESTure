package bosh

import (
	"github.com/evoila/infraTESTure/config"
	"github.com/fatih/color"
	"log"
	"strconv"
)

type Bosh struct {}

var conf *config.Config

// Initialize the Bosh Director and the Deployment affiliated to the deployment name in the config
func InitInfrastructureValues(config *config.Config) {
	conf = config

	deploymentName = config.DeploymentName

	var err error

	// Create a bosh director to get the deployment
	boshDirector, err = buildDirector(config)
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
		log.Printf(color.RedString("[ERROR] " + customMessage + ": " + err.Error()))
	} else {
		log.Printf(color.RedString("[ERROR] " + customMessage))
	}
}

func stringToFloat(value string) float64 {
	floatValue, err := strconv.ParseFloat(value, 64)

	if err != nil {
		logError(err, "")
	}

	return floatValue
}
