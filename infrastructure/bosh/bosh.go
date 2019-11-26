package bosh

import (
	"bytes"
	"fmt"
	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/evoila/infraTESTure/config"
	"github.com/evoila/infraTESTure/infrastructure"
	"github.com/fatih/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Bosh struct {}

const (
	netem =	"sudo tc qdisc add dev eth0 root handle 1a1a: htb default 1 && " +
			"sudo tc class add dev eth0 parent 1a1a: classid 1a1a:1 htb rate 10000000.0kbit && " +
			"sudo tc class add dev eth0 parent 1a1a: classid 1a1a:2 htb rate 10000000.0Kbit ceil 10000000.0Kbit && " +
			"sudo tc filter add dev eth0 protocol ip parent 1a1a: prio 1 u32 match ip sport 22 0xffff flowid 1a1a:1 && " +
			"sudo tc filter add dev eth0 protocol ip parent 1a1a: prio 1 u32 match ip dst %s flowid 1a1a:1 && " +
			"sudo tc filter add dev eth0 protocol ip parent 1a1a: prio 2 u32 match ip src 0.0.0.0/0 match ip dst 0.0.0.0/0 flowid 1a1a:2 && " +
			"sudo tc qdisc add dev eth0 parent 1a1a:2 handle 2518: netem %s"

	removeTC = "sudo tc qdisc del dev eth0 root"
	cpuLoad = "sudo apt-get -y install stress-ng && setsid stress-ng -c 1 -l %d &>/dev/null"
	memLoad = "sudo apt-get -y install stress-ng && setsid stress-ng --vm-bytes $(awk '/MemAvailable/{printf \"%%d\\n\", $2 * %f;}' < /proc/meminfo)k -m 1 &>/dev/null"
	stopStress = "sudo kill $(pgrep -o -x stress-ng) && sudo apt-get -y remove stress-ng"
)

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

func (b *Bosh) FillDisk(size int, path string, fileName string, vmId string) {
	RunSshCommand(vmId, fmt.Sprintf("cd %s && sudo dd if=/dev/urandom of=%s count=%v bs=1048576", path, fileName, size))
}

func (b *Bosh) CleanupDisk(path string, fileName string, vmId string) {
	_, err := RunSshCommand(vmId, fmt.Sprintf("sudo rm %s/%s", path, fileName))

	if err != nil {
		logError(err, "Cleanup disk failed")
	}
}

func (b *Bosh) SimulatePackageLoss(loss int, correlation int) string {
	if loss < 0 || loss > 100 || correlation < 0 || correlation > 100 {
		logError(nil, "Invalid value. Loss and correlation cannot be lower than 0 or greater than 100")
		return ""
	}

	return fmt.Sprintf("loss %d%% %d%%", loss, correlation)
}

func (b *Bosh) SimulatePackageCorruption(corruption int, correlation int) string {
	if corruption < 0 || corruption > 100 || correlation < 0 || correlation > 100 {
		logError(nil,"Invalid value. Corruption and correlation cannot be lower than 0 or greater than 100")
		return ""
	}

	return fmt.Sprintf("corrupt %d%% %d%%", corruption, correlation)
}

func (b *Bosh) SimulatePackageDuplication(duplication int, correlation int) string {
	if duplication < 0 || duplication > 100 || correlation < 0 || correlation > 100 {
		logError(nil, "Invalid value. Duplication and correlation cannot be lower than 0 or greater than 100")
		return ""
	}

	return fmt.Sprintf("duplicate %d%% %d%%", duplication, correlation)
}

func (b *Bosh) SimulateNetworkDelay(delay int, variation int) string {
	if delay < 0 || variation < 0 {
		logError(nil, "Invalid value. Delay and variation cannot be lower than 0")
		return ""
	}

	if variation > 0 {
		return fmt.Sprintf("delay %dms %dms distribution normal", delay, variation)
	} else {
		return fmt.Sprintf("delay %dms", delay)
	}
}

func (b *Bosh) AddTrafficControl(vmId string, directorIp string, tc string) {
	_, err := RunSshCommand(vmId, fmt.Sprintf(netem, directorIp, tc))

	if err != nil {
		logError(err, "Failed to simulate traffic control")
	}
}

func (b *Bosh) RemoveTrafficControl(vmId string) {
	session, client, err := createSshSession(vmId)
	defer client.Close()
	defer session.Close()

	if err != nil {
		logError(err, "Failed to create SSH session")
	}

	_, err = RunSshCommand(vmId, removeTC)

	if err != nil {
		logError(err, "Failed to remove Traffic Control")
	}
}

func (b *Bosh) StartCPULoad(vmId string, percentage int) {
	_, err := RunSshCommand(vmId, fmt.Sprintf(cpuLoad, percentage))

	if err != nil {
		logError(err, "Failed to simulate CPU load")
	}
}

func (b *Bosh) StartMemLoad(vmId string, percentage float64) {
	_, err := RunSshCommand(vmId, fmt.Sprintf(memLoad, percentage/100))

	if err != nil {
		logError(err, "Failed to simulate Memory Load")
	}
}

func (b *Bosh) StopStress(vmId string) {
	_, err := RunSshCommand(vmId, fmt.Sprintf(stopStress))

	if err != nil {
		logError(err, "Failed to kill stress process")
	}
}

func RunSshCommand(vmId string, command string) (string, error) {
	session, client, err := createSshSession(vmId)
	defer client.Close()
	defer session.Close()

	if err != nil {
		logError(err, "Failed to create SSH session")
		return "", err
	}

	var result bytes.Buffer
	session.Stdout = &result
	session.Stderr = os.Stderr

	err = session.Run(command)

	return result.String(), err
}

func logError(err error, customMessage string) {
	if err != nil {
		log.Fatal(color.RedString("[ERROR] " + customMessage + ": " + err.Error()))
	} else {
		log.Fatal(color.RedString("[ERROR] " + customMessage))
	}
}

func stringToFloat(value string) float64 {
	floatValue, err := strconv.ParseFloat(value, 64)

	if err != nil {
		logError(err, "")
	}

	return floatValue
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