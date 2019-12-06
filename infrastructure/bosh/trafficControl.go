package bosh

import "fmt"

const (
	netem =	"sudo tc qdisc add dev eth0 root handle 1a1a: htb default 1 && " +
		"sudo tc class add dev eth0 parent 1a1a: classid 1a1a:1 htb rate 10000000.0kbit && " +
		"sudo tc class add dev eth0 parent 1a1a: classid 1a1a:2 htb rate 10000000.0Kbit ceil 10000000.0Kbit && " +
		"sudo tc filter add dev eth0 protocol ip parent 1a1a: prio 1 u32 match ip sport 22 0xffff flowid 1a1a:1 && " +
		"sudo tc filter add dev eth0 protocol ip parent 1a1a: prio 1 u32 match ip dst %s flowid 1a1a:1 && " +
		"sudo tc filter add dev eth0 protocol ip parent 1a1a: prio 2 u32 match ip src 0.0.0.0/0 match ip dst 0.0.0.0/0 flowid 1a1a:2 && " +
		"sudo tc qdisc add dev eth0 parent 1a1a:2 handle 2518: netem %s"

	removeTC = "sudo tc qdisc del dev eth0 root"
)

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