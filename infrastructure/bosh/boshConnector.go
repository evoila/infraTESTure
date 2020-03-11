package bosh

import (
	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-cli/uaa"
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/evoila/infraTESTure/config"
	"io/ioutil"
)

var boshDirector director.Director
var deployment director.Deployment
var deploymentName string

// Build a Director based on the director URL from the configuration
// @param config Initialized config struct from github.com/evoila/infraTESTure/config
// @return director Initialized director struct from github.com/cloudfoundry/bosh-cli/director
func BuildDirector(config *config.Config) (director.Director, error) {
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
	// If a ca file is provided use the files content instead of the cert from the yaml
	factoryConfig.CACert = setCa(config)

	// Allow Director to fetch UAA tokens when necessary.
	boshUaa, err := BuildUAA(config)
	if err != nil {
		return nil, err
	}
	factoryConfig.TokenFunc = uaa.NewClientTokenSession(boshUaa).TokenFunc

	return factory.New(factoryConfig, director.NewNoopTaskReporter(), director.NewNoopFileReporter())
}

// Build an UAA based on the UAA URL from the configuration
// @param config Initialized config struct from github.com/evoila/infraTESTure/config
// @return uaa Initialized uaa struct from github.com/cloudfoundry/bosh-cli/uaa
func BuildUAA(config *config.Config) (uaa.UAA, error) {
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

	// Configure custom trusted CA certificates.
	// If nothing is provided default system certificates are used.
	// If a ca file is provided use the files content instead of the cert from the yaml
	boshConfig.CACert = setCa(config)

	return factory.New(boshConfig)
}

func setCa(config *config.Config) string {
	var ca = config.Bosh.Ca

	if conf.Bosh.CaFile != "" {
		ca = readCaFromFile(config.Bosh.CaFile)
	}

	return ca
}

func readCaFromFile(pathToCert string) string {
	content, err := ioutil.ReadFile(pathToCert)
	if err != nil {
		panic(err)
	}

	return string(content)
}
