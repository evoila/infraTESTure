package bosh

import (
	"fmt"

	"github.com/evoila/infraTESTure/infrastructure"
)

const (
	netem = "IFACE=`sudo netstat -ie | grep -B1 \"%s\" | head -n1 | awk '{print $1}'` &&" +
		"sudo tc qdisc add dev $IFACE root handle 1a1a: htb default 1 && " +
		"sudo tc class add dev $IFACE parent 1a1a: classid 1a1a:1 htb rate 10000000.0kbit && " +
		"sudo tc class add dev $IFACE parent 1a1a: classid 1a1a:2 htb rate 10000000.0Kbit ceil 10000000.0Kbit && " +
		"sudo tc filter add dev $IFACE protocol ip parent 1a1a: prio 1 u32 match ip sport 22 0xffff flowid 1a1a:1 && " +
		"sudo tc filter add dev $IFACE protocol ip parent 1a1a: prio 1 u32 match ip dst %s flowid 1a1a:1 && " +
		"sudo tc filter add dev $IFACE protocol ip parent 1a1a: prio 2 u32 match ip src 0.0.0.0/0 match ip dst 0.0.0.0/0 flowid 1a1a:2 && " +
		"sudo tc qdisc add dev $IFACE parent 1a1a:2 handle 2518: netem %s"

	removeTC = "IFACE=`sudo netstat -ie | grep -B1 \"%s\" | head -n1 | awk '{print $1}'` && sudo tc qdisc del dev $IFACE root"
)

// Create a big dump file with a given size in MB
func (b *Bosh) FillDisk(size int, path string, fileName string, vmId string) {
	RunSshCommand(vmId, fmt.Sprintf("cd %s && sudo dd if=/dev/urandom of=%s count=%v bs=1048576", path, fileName, size))
}

// Remove dump file
func (b *Bosh) CleanupDisk(path string, fileName string, vmId string) {
	_, err := RunSshCommand(vmId, fmt.Sprintf("sudo rm %s/%s", path, fileName))

	if err != nil {
		logError(err, "Cleanup disk failed")
	}
}

// The creation of the tc commands and the actual execution are separated here
// because the user may wants to create a more complex tc command consisting of
// e.g. package loss AND network delay

// Creates tc command for package loss based on given parameters
func (b *Bosh) SimulatePackageLoss(loss int, correlation int) string {
	if loss < 0 || loss > 100 || correlation < 0 || correlation > 100 {
		logError(nil, "Invalid value. Loss and correlation cannot be lower than 0 or greater than 100")
		return ""
	}

	return fmt.Sprintf("loss %d%% %d%%", loss, correlation)
}

// Creates tc command for package corruption based on given parameters
func (b *Bosh) SimulatePackageCorruption(corruption int, correlation int) string {
	if corruption < 0 || corruption > 100 || correlation < 0 || correlation > 100 {
		logError(nil, "Invalid value. Corruption and correlation cannot be lower than 0 or greater than 100")
		return ""
	}

	return fmt.Sprintf("corrupt %d%% %d%%", corruption, correlation)
}

// Creates tc command for package duplication based on given parameters
func (b *Bosh) SimulatePackageDuplication(duplication int, correlation int) string {
	if duplication < 0 || duplication > 100 || correlation < 0 || correlation > 100 {
		logError(nil, "Invalid value. Duplication and correlation cannot be lower than 0 or greater than 100")
		return ""
	}

	return fmt.Sprintf("duplicate %d%% %d%%", duplication, correlation)
}

// Creates tc command for network delay based on given parameters
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

// Open SSH session and run the tc command
func (b *Bosh) AddTrafficControl(vm infrastructure.VM, directorIp string, tc string) {
	command := fmt.Sprintf(netem, vm.IPs[0], directorIp, tc)
	_, err := RunSshCommand(vm.ID, command)

	if err != nil {
		logError(err, "Failed to simulate traffic control")
	}
}

// Open SSH session and run command for removing the traffic control
func (b *Bosh) RemoveTrafficControl(vm infrastructure.VM) {
	_, err := RunSshCommand(vm.ID, fmt.Sprintf(removeTC, vm.IPs[0]))

	if err != nil {
		logError(err, "Failed to remove Traffic Control")
	}
}
