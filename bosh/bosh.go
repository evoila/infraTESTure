package bosh

import (
	"github.com/briandowns/spinner"
	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-cli/uaa"
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/evoila/infraTESTure/config"
	"log"
	"strconv"
	"time"
)

var spin *spinner.Spinner
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


func InitInfrastructureValues(config *config.Config) {
	spin = spinner.New(spinner.CharSets[33], 100*time.Millisecond)

	deploymentName = config.DeploymentName

	// Create a bosh director to get the deployment
	boshDirector, err := buildDirector(config)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	// See director/interfaces.go for a full list of methods.
	deployment, err = boshDirector.FindDeployment(config.DeploymentName)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}
}

// Check if all VMs of a deployment are running
func IsDeploymentRunning() bool {

	log.Printf("[INFO] Checking process state for every VM of Deployment %v...", deploymentName)

	spin.Start()

	// Get all vitals information about the deployment
	vmVitals, err := deployment.VMInfos()

	if err != nil {
		spin.Stop()
		log.Printf("[ERROR] %v", err)
	}

	healthy := true

	spin.Stop()

	// Check if one of the VMs has a different process state than "running"
	for _, vmVital := range vmVitals {
		log.Printf("[INFO] %v/%v - %v", vmVital.JobName, vmVital.ID, vmVital.ProcessState)

		if !vmVital.IsRunning() {
			healthy = false
		}
	}

	return healthy
}

// Stop one VM of the deployment
// TODO: Stop a random number of VMs, but not all
func Stop(index int) {

	log.Printf("[INFO] Shutting down vm...")

	spin.Start()

	err := deployment.Stop(director.NewAllOrInstanceGroupOrInstanceSlug("redis", strconv.Itoa(index)), director.StopOpts{})

	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	spin.Stop()
}

// Restart a stopped VM
func Start(index int) {

	log.Printf("[INFO] Restarting vm...")

	spin.Start()

	err := deployment.Start(director.NewAllOrInstanceGroupOrInstanceSlug("redis", strconv.Itoa(index)), director.StartOpts{})

	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	spin.Stop()
}

// Get the ips of all VMs of a deployment
func GetIPs() []string {
	vmInfos, err := deployment.VMInfos()

	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	var ips []string

	for _, vmInfo := range vmInfos {
		for _, ip := range vmInfo.IPs {
			ips = append(ips, ip)
		}
	}

	return ips
}