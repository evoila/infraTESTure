package infrastructure

type Infrastructure interface {
	Start(int)
	Stop(int)
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
	MemoryUsage float64
	DiskUsage float64
}