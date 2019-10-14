package bosh

import (
	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/evoila/infraTESTure/config"
	"github.com/evoila/infraTESTure/infrastructure"
	"github.com/fatih/color"
	"log"
	"math"
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

// Stop a VM based on the VM ID
func (b *Bosh) Stop(id string) {
	err := deployment.Stop(director.NewAllOrInstanceGroupOrInstanceSlug("", id), director.StopOpts{
		Converge:    true,
	})

	if err != nil {
		logError(err, "")
	}
}

// Start a VM based on the VM ID
func (b *Bosh) Start(id string) {
	// Restart a stopped VM
	err := deployment.Start(director.NewAllOrInstanceGroupOrInstanceSlug("", id), director.StartOpts{
		Converge:    true,
	})

	if err != nil {
		logError(err, "")
	}
}

// Return a map of all IPs of the deployment with the VM ID as the key, and all affiliated IPs as the value
func (b *Bosh) GetIPs() map[string][]string {
	// Get the ips of all VMs of a deployment

	vmInfos, err := deployment.VMInfos()

	if err != nil {
		logError(err, "")
	}

	ips := make(map[string][]string)

	for _, vmInfo := range vmInfos {
		for _, ip := range vmInfo.IPs {
				ips[vmInfo.ID] = append(ips[vmInfo.ID], ip)
		}
	}

	return ips
}

// Return an own Deployment struct with some important metrics
func (b *Bosh) GetDeployment() infrastructure.Deployment {
	vmVitals, err := deployment.VMInfos()

	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	var vms []infrastructure.VM

	for _, vmVital := range vmVitals {
		vms = append(vms, infrastructure.VM{
			ServiceName:           vmVital.JobName,
			ID:                    vmVital.ID,
			IPs:				   vmVital.IPs,
			State:                 vmVital.ProcessState,
			CpuUsage:          	   math.Round(stringToFloat(vmVital.Vitals.CPU.User) + stringToFloat(vmVital.Vitals.CPU.Sys)),
			MemoryUsagePercentage: stringToFloat(vmVital.Vitals.Mem.Percent),
			MemoryUsageTotal:      stringToFloat(vmVital.Vitals.Mem.KB),
			//TODO: find out how to get disk size
			DiskSize:			   0,
			DiskUsage:			   stringToFloat(vmVital.Vitals.Disk["system"].Percent) +
				                   stringToFloat(vmVital.Vitals.Disk["ephemeral"].Percent) +
								   stringToFloat(vmVital.Vitals.Disk["persistent"].Percent),
		})
	}

	return infrastructure.Deployment{
		DeploymentName: deploymentName,
		Hosts:          b.GetIPs(),
		VMs:            vms,
	}
}

// Check if a VM is running
func (b *Bosh) IsRunning() bool {
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

func stringToFloat(value string) float64 {
	floatValue, err := strconv.ParseFloat(value, 64)

	if err != nil {
		logError(err, "")
	}

	return floatValue
}