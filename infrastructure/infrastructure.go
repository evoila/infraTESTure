package infrastructure

// Generic Infrastructure interface used by the actual test repository. This will be
// initialized at the runtime
type Infrastructure interface {
	Start(string)
	Stop(string)
	GetIPs() map[string][]string
	GetDeployment() Deployment
	IsRunning() bool
	FillDisk(int, string, string, string)
	CleanupDisk(string, string, string)
}

type Deployment struct {
	DeploymentName string
	Hosts map[string][]string
	VMs []VM
}

type VM struct {
	ServiceName string
	ID string
	IPs []string
	State string
	DiskSize float64
	DiskUsageTotal float64
	DiskUsagePercentage float64
	CpuUsage float64
	MemoryUsagePercentage float64
	MemoryUsageTotal float64
}