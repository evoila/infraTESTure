package infrastructure

type Infrastructure interface {
	Start(string)
	Stop(string)
	GetIPs() []string
	GetDeployment() Deployment
	IsRunning() bool
}

type Deployment struct {
	DeploymentName string
	Hosts []string
	VMs []VM
}

type VM struct {
	ServiceName string
	ID string
	State string
	DiskSize float64
	CpuUsage float64
	MemoryUsagePercentage float64
	MemoryUsageTotal float64
	DiskUsage float64
}