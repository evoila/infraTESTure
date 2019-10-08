package bosh

import (
	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-cli/uaa"
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/evoila/infraTESTure/config"
)

var deployment director.Deployment
var deploymentName string

func buildDirector(config *config.Config) (director.Director, error) {
	logs := logger.NewLogger(logger.LevelError)
	factory := director.NewFactory(logs)

	// Build a Director config from address-like string.
	// HTTPS is required and certificates are always verified.
	factoryConfig, err := director.NewConfigFromURL(config.Bosh.DirectorUrl)
	if err != nil {
		return nil, err
	}

	// Configure custom trusted CA certificates.
	// If nothing is provided default system certificates are used.
	factoryConfig.CACert = config.Bosh.Ca

	// Allow Director to fetch UAA tokens when necessary.
	boshUaa, err := buildUAA(config)
	if err != nil {
		return nil, err
	}
	factoryConfig.TokenFunc = uaa.NewClientTokenSession(boshUaa).TokenFunc

	return factory.New(factoryConfig, director.NewNoopTaskReporter(), director.NewNoopFileReporter())
}

func buildUAA(config *config.Config) (uaa.UAA, error) {
	logs := logger.NewLogger(logger.LevelError)
	factory := uaa.NewFactory(logs)

	// Build a UAA config from a URL.
	// HTTPS is required and certificates are always verified.
	boshConfig, err := uaa.NewConfigFromURL(config.Bosh.UaaUrl)
	if err != nil {
		return nil, err
	}

	// Set client credentials for authentication.
	// Machine level access should typically use a client instead of a particular user.
	boshConfig.Client = config.Bosh.UaaClient
	boshConfig.ClientSecret = config.Bosh.UaaClientSecret

	// Configure trusted CA certificates.
	// If nothing is provided default system certificates are used.
	boshConfig.CACert = config.Bosh.Ca

	return factory.New(boshConfig)
}