package bosh

import (
	"bytes"
	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/evoila/infraTESTure/infrastructure"
	"log"
	"math"
	"os"
	"strings"
)

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
}x^

// Return an own Deployment struct with some important metrics
func (b *Bosh) GetDeployment() infrastructure.Deployment {
	vmVitals, err := deployment.VMInfos()

	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	var vms []infrastructure.VM

	for _, vmVital := range vmVitals {
		used, available := ParseDiskSize(vmVital.ID)

		vms = append(vms, infrastructure.VM{
			ServiceName:           vmVital.JobName,
			ID:                    vmVital.ID,
			IPs:				   vmVital.IPs,
			State:                 vmVital.ProcessState,
			CpuUsage:          	   math.Round(stringToFloat(vmVital.Vitals.CPU.User) + stringToFloat(vmVital.Vitals.CPU.Sys)),
			MemoryUsagePercentage: stringToFloat(vmVital.Vitals.Mem.Percent),
			MemoryUsageTotal:      stringToFloat(vmVital.Vitals.Mem.KB),
			DiskSize:			   stringToFloat(used) + stringToFloat(available),
			DiskUsageTotal:		   stringToFloat(used),
			DiskUsagePercentage:   stringToFloat(vmVital.Vitals.Disk["persistent"].Percent),
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

func ParseDiskSize(vmId string) (used string, available string){
	session, client, err := createSshSession(vmId)
	defer client.Close()
	defer session.Close()

	if err != nil && client == nil {
		logError(err, "Failed to create SSH session")
	}

	var result bytes.Buffer
	session.Stdout = &result
	session.Stderr = os.Stderr

	err = session.Run("df | awk 'NR > 1{print $6\" \"$3\" \"$4 }'")

	fields := strings.Fields(result.String())

	for i, field := range fields {
		//TODO: Hardcoded or configurable?
		if strings.HasPrefix(field, "/var/vcap/store") {
			return fields[i+1], fields[i+2]
		}
	}

	return "", ""
}