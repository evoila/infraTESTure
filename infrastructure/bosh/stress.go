package bosh

import "fmt"

const (
	cpuLoad = "sudo apt-get -y install stress-ng && setsid stress-ng -c 1 -l %d &>/dev/null"
	memLoad = "sudo apt-get -y install stress-ng && setsid stress-ng --vm-bytes $(awk '/MemAvailable/{printf \"%%d\\n\", $2 * %f;}' < /proc/meminfo)k -m 1 &>/dev/null"
	stopStress = "sudo kill $(pgrep -o -x stress-ng) && sudo apt-get -y remove stress-ng"
)

// SSH to vm, install stress-ng and increase CPU load by a given percentage
// @param vmId Id of the vm you want to increase the CPU load
// @param percentage Percentage value you want to increase the CPU load to
func (b *Bosh) StartCPULoad(vmId string, percentage int) {
	_, err := RunSshCommand(vmId, fmt.Sprintf(cpuLoad, percentage))

	if err != nil {
		logError(err, "Failed to simulate CPU load")
	}
}

// SSH to vm, install stress-ng and increase RAM load by a given percentage
// @param vmId Id of the vm you want to increase the Memory load
// @param percentage Percentage value you want to increase the Memory load to
func (b *Bosh) StartMemLoad(vmId string, percentage float64) {
	_, err := RunSshCommand(vmId, fmt.Sprintf(memLoad, percentage/100))

	if err != nil {
		logError(err, "Failed to simulate Memory Load")
	}
}

// SSH to vm, kill the stress process and uninstall stress-ng
// @param vmId Id of the VM you want to kill the stress process on
func (b *Bosh) StopStress(vmId string) {
	_, err := RunSshCommand(vmId, fmt.Sprintf(stopStress))

	if err != nil {
		logError(err, "Failed to kill stress process")
	}
}