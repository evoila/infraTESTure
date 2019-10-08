package bosh

import (
	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/evoila/infraTESTure/config"
	"github.com/evoila/infraTESTure/infrastructure"
	"github.com/fatih/color"
	"log"
	"strconv"
)

type Bosh struct {}

func InitInfrastructureValues(config *config.Config) {
	deploymentName = config.DeploymentName

	// Create a bosh director to get the deployment
	boshDirector, err := buildDirector(config)
	if err != nil {
		logError(err, "")
	}

	// See director/interfaces.go for a full list of methods.
	deployment, err = boshDirector.FindDeployment(config.DeploymentName)
	if err != nil {
		logError(err, "")
	}
}

func (b Bosh) Stop(index int) {
	// Stop one VM of the deployment
	log.Printf("[INFO] Shutting down vm...")

	err := deployment.Stop(director.NewAllOrInstanceGroupOrInstanceSlug("redis", strconv.Itoa(index)), director.StopOpts{})

	if err != nil {
		logError(err, "")
	}
}

func (b Bosh) Start(index int) {
	// Restart a stopped VM

	log.Printf("[INFO] Restarting vm...")

	err := deployment.Start(director.NewAllOrInstanceGroupOrInstanceSlug("redis", strconv.Itoa(index)), director.StartOpts{})

	if err != nil {
		logError(err, "")
	}
}

func (b Bosh) GetIPs() []string {
	// Get the ips of all VMs of a deployment

	vmInfos, err := deployment.VMInfos()

	if err != nil {
		logError(err, "")
	}

	var ips []string

	for _, vmInfo := range vmInfos {
		for _, ip := range vmInfo.IPs {
			ips = append(ips, ip)
		}
	}

	return ips
}

func (b Bosh) GetDeployment() infrastructure.Deployment {
	vmVitals, err := deployment.VMInfos()

	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	var vms []infrastructure.VM

	// Check if one of the VMs has a different process state than "running"
	for _, vmVital := range vmVitals {
		vms = append(vms, infrastructure.VM{
			ServiceName: vmVital.JobName,
			ID:			 vmVital.ID,
			State:       vmVital.ProcessState,
			DiskSize:    0, //TODO
			CpuUsage:    0, //TODO
			MemoryUsage: 0, //TODO
			DiskUsage:   0, //TODO
		})
	}

	return infrastructure.Deployment{
		DeploymentName: deploymentName,
		Hosts:          b.GetIPs(),
		VMs:            vms,
	}
}

func (b Bosh) IsRunning() bool {
	vmVitals, err := deployment.VMInfos()

	if err != nil {
		log.Printf("[ERROR] %v\n", err)
	}

	for _, vmVital := range vmVitals {
		if vmVital.ProcessState != "running" {
			return false
		}
	}

	return true
}

func logError(err error, customMessage string) {
	log.Fatal(color.RedString("[ERROR] " + customMessage + ": " + err.Error()))
}