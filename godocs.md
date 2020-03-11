# bosh
--
    import "github.com/evoila/infraTESTure/infrastructure/bosh"


## Usage

#### func  InitInfrastructureValues

```go
func InitInfrastructureValues(config *config.Config)
```
Initialize the Bosh Director and the Deployment affiliated to the deployment
name in the config

#### func  ParseDiskSize

```go
func ParseDiskSize(vmId string) (used string, available string)
```
SSH to vm and run df command in order to get the free disk space then filter the
one needed

#### func  RunSshCommand

```go
func RunSshCommand(vmId string, command string) (string, error)
```

#### type Bosh

```go
type Bosh struct{}
```


#### func (*Bosh) AddTrafficControl

```go
func (b *Bosh) AddTrafficControl(vmId string, directorIp string, tc string)
```
Open SSH session and run the tc command

#### func (*Bosh) AssertEquals

```go
func (b *Bosh) AssertEquals(actual interface{}, expected interface{}) bool
```

#### func (*Bosh) AssertFalse

```go
func (b *Bosh) AssertFalse(value bool) bool
```

#### func (*Bosh) AssertNil

```go
func (b *Bosh) AssertNil(value interface{}) bool
```

#### func (*Bosh) AssertNotEquals

```go
func (b *Bosh) AssertNotEquals(actual interface{}, expected interface{}) bool
```

#### func (*Bosh) AssertNotNil

```go
func (b *Bosh) AssertNotNil(value interface{}) bool
```

#### func (*Bosh) AssertTrue

```go
func (b *Bosh) AssertTrue(value bool) bool
```

#### func (*Bosh) CleanupDisk

```go
func (b *Bosh) CleanupDisk(path string, fileName string, vmId string)
```
Remove dump file

#### func (*Bosh) FillDisk

```go
func (b *Bosh) FillDisk(size int, path string, fileName string, vmId string)
```
Create a big dump file with a given size in MB

#### func (*Bosh) GetDeployment

```go
func (b *Bosh) GetDeployment() infrastructure.Deployment
```
Return an own Deployment struct with some important metrics

#### func (*Bosh) GetIPs

```go
func (b *Bosh) GetIPs() map[string][]string
```
Return a map of all IPs of the deployment with the VM ID as the key, and all
affiliated IPs as the value

#### func (*Bosh) IsRunning

```go
func (b *Bosh) IsRunning() bool
```
Check if a VM is running

#### func (*Bosh) RemoveTrafficControl

```go
func (b *Bosh) RemoveTrafficControl(vmId string)
```
Open SSH session and run command for removing the traffic control

#### func (*Bosh) SimulateNetworkDelay

```go
func (b *Bosh) SimulateNetworkDelay(delay int, variation int) string
```
Creates tc command for network delay based on given parameters

#### func (*Bosh) SimulatePackageCorruption

```go
func (b *Bosh) SimulatePackageCorruption(corruption int, correlation int) string
```
Creates tc command for package corruption based on given parameters

#### func (*Bosh) SimulatePackageDuplication

```go
func (b *Bosh) SimulatePackageDuplication(duplication int, correlation int) string
```
Creates tc command for package duplication based on given parameters

#### func (*Bosh) SimulatePackageLoss

```go
func (b *Bosh) SimulatePackageLoss(loss int, correlation int) string
```
Creates tc command for package loss based on given parameters

#### func (*Bosh) Start

```go
func (b *Bosh) Start(id string)
```
Start a VM based on the VM ID

#### func (*Bosh) StartCPULoad

```go
func (b *Bosh) StartCPULoad(vmId string, percentage int)
```
SSH to vm, install stress-ng and increase CPU load by a given percentage

#### func (*Bosh) StartMemLoad

```go
func (b *Bosh) StartMemLoad(vmId string, percentage float64)
```
SSH to vm, install stress-ng and increase RAM load by a given percentage

#### func (*Bosh) Stop

```go
func (b *Bosh) Stop(id string)
```
Stop a VM based on the VM ID

#### func (*Bosh) StopStress

```go
func (b *Bosh) StopStress(vmId string)
```
SSH to vm, kill the stress process and uninstall stress-ng
