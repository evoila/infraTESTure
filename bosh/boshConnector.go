package bosh

import (
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshuaa "github.com/cloudfoundry/bosh-cli/uaa"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/evoila/infraTESTure/config"
)

func buildDirector(config *config.Config) (boshdir.Director, error) {
	logger := boshlog.NewLogger(boshlog.LevelError)
	factory := boshdir.NewFactory(logger)

	// Build a Director config from address-like string.
	// HTTPS is required and certificates are always verified.
	factoryConfig, err := boshdir.NewConfigFromURL(config.Bosh.DirectorUrl)
	if err != nil {
		return nil, err
	}

	// Configure custom trusted CA certificates.
	// If nothing is provided default system certificates are used.
	factoryConfig.CACert = config.Bosh.Ca

	// Allow Director to fetch UAA tokens when necessary.
	uaa, err := buildUAA(config)
	if err != nil {
		return nil, err
	}
	factoryConfig.TokenFunc = boshuaa.NewClientTokenSession(uaa).TokenFunc

	return factory.New(factoryConfig, boshdir.NewNoopTaskReporter(), boshdir.NewNoopFileReporter())
}

func buildUAA(config *config.Config) (boshuaa.UAA, error) {
	logger := boshlog.NewLogger(boshlog.LevelError)
	factory := boshuaa.NewFactory(logger)

	// Build a UAA config from a URL.
	// HTTPS is required and certificates are always verified.
	boshConfig, err := boshuaa.NewConfigFromURL(config.Bosh.UaaUrl)
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
