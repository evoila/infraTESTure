# bosh
import "github.com/evoila/infraTESTure/infrastructure/bosh"

## Usage

#### func  BuildDirector

```go
func BuildDirector(config *config.Config) (director.Director, error)
```
Build a Director based on the director URL from the configuration 

@param config Initialized config struct from github.com/evoila/infraTESTure/config 

@return director Initialized director struct from github.com/cloudfoundry/bosh-cli/director

#### func  BuildUAA

```go
func BuildUAA(config *config.Config) (uaa.UAA, error)
```
Build an UAA based on the UAA URL from the configuration 

@param config Initialized config struct from github.com/evoila/infraTESTure/config 

@return uaa Initialized uaa struct from github.com/cloudfoundry/bosh-cli/uaa

#### func  InitInfrastructureValues

```go
func InitInfrastructureValues(config *config.Config)
```
Initialize the Bosh Director and the Deployment affiliated to the deployment name in the config 

@param config Configuration struct from github.com/evoila/infraTESTure/config

#### func  ParseDiskSize

```go
func ParseDiskSize(vmId string) (used string, available string)
```
SSH to vm and run df command in order to get the free disk space then filter the one needed 

@param vmId Id of the VM you want to determine used and available disk space from @return used Value in bytes of used disk space 

@return available Value in bytes of available disk space

#### func  RunSshCommand

```go
func RunSshCommand(vmId string, command string) (string, error)
```
Run a ssh command on a VM @param vmId Id of the VM you want to run the command on 

@param command Command you want to execute on the VM 

@return string Stdout of the command execution

#### type Bosh

```go
type Bosh struct{}
```


#### func (*Bosh) AddTrafficControl

```go
func (b *Bosh) AddTrafficControl(vmId string, directorIp string, tc string)
```
Open SSH session and run the tc command 

@param vmId Id of the VM you want to add traffic control to 

@param tc String containing the tc command, from one or more of the previous functions.

TC commands can be connected. For example you can run the "SimulatePackageLoss" function and afterwards expand the resulting string with "SimulateNetworkDelay"

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

@param path Path on the VM where the dump file was created

@param fileName Name of the dump file @param vmId Id of the VM you want to clean the persistent disk again

#### func (*Bosh) FillDisk

```go
func (b *Bosh) FillDisk(size int, path string, fileName string, vmId string)
```
Create a big dump file with a given size in MB @param size Size of the file that is being created on the VM 

@param path Path on the VM where the dump file is saved to 

@param fileName Name of the dump file @param vmId Id of the VM you want to fill the persistent disk on

#### func (*Bosh) GetDeployment

```go
func (b *Bosh) GetDeployment() infrastructure.Deployment
```
Return an own Deployment struct with some important metrics 

@return deployment Initialized deployment struct from github.com/evoila/infraTESTure/infrastructure

#### func (*Bosh) GetIPs

```go
func (b *Bosh) GetIPs() map[string][]string
```
Return a map of all IPs of the deployment with the VM ID as the key, and all affiliated IPs as the value 

@return map Map containing all IPs of the deployment

#### func (*Bosh) IsRunning

```go
func (b *Bosh) IsRunning() bool
```
Check if a VM is running 

@return bool Bool value telling if the VM is running (true) or not (false)

#### func (*Bosh) RemoveTrafficControl

```go
func (b *Bosh) RemoveTrafficControl(vmId string)
```
Open SSH session and run command for removing the traffic control 

@param vmId Id of the VM you want to remove the traffic control from

#### func (*Bosh) SimulateNetworkDelay

```go
func (b *Bosh) SimulateNetworkDelay(delay int, variation int) string
```
Creates tc command for network delay based on given parameters 

@param delay Value in ms of how the network communication should be delayed 

@param variation Optional value to variate the delay value 

@return string String containing a TC command that can be used as a parameter for the AddTrafficControl function

#### func (*Bosh) SimulatePackageCorruption

```go
func (b *Bosh) SimulatePackageCorruption(corruption int, correlation int) string
```
Creates tc command for package corruption based on given parameters 

@param corruption Percentage value of the package corruption that should be simulated

@param correlation Optional correlation value to decide where a package should be corrupted or not, based on the decision of the previous package 

@return string String containing a TC command that can be used as a parameter for the
AddTrafficControl function

#### func (*Bosh) SimulatePackageDuplication

```go
func (b *Bosh) SimulatePackageDuplication(duplication int, correlation int) string
```
Creates tc command for package duplication based on given parameters 

@param duplication Percentage value of the package duplication that should be simulated

@param correlation Optional correlation value to decide where a package should
be duplicated or not, based on the decision of the previous package 

@return string String containing a TC command that can be used as a parameter for the AddTrafficControl function

#### func (*Bosh) SimulatePackageLoss

```go
func (b *Bosh) SimulatePackageLoss(loss int, correlation int) string
```
Creates tc command for package loss based on given parameters 

@param loss Percentage value of the package loss that should be simulated 

@param correlation Optional correlation value to decide where a package should be dropped or not, based on the decision of the previous package 

@return string String containing a TC command that can be used as a parameter for the AddTrafficControl function

#### func (*Bosh) Start

```go
func (b *Bosh) Start(id string)
```
Start a VM based on the VM ID 

@param id Id of the VM you want to start

#### func (*Bosh) StartCPULoad

```go
func (b *Bosh) StartCPULoad(vmId string, percentage int)
```
SSH to vm, install stress-ng and increase CPU load by a given percentage 

@param vmId Id of the vm you want to increase the CPU load 

@param percentage Percentage value you want to increase the CPU load to

#### func (*Bosh) StartMemLoad

```go
func (b *Bosh) StartMemLoad(vmId string, percentage float64)
```
SSH to vm, install stress-ng and increase RAM load by a given percentage 

@param vmId Id of the vm you want to increase the Memory load 

@param percentage Percentage value you want to increase the Memory load to

#### func (*Bosh) Stop

```go
func (b *Bosh) Stop(id string)
```
Stop a VM based on the VM ID 

@param id Id of the VM you want to stop

#### func (*Bosh) StopStress

```go
func (b *Bosh) StopStress(vmId string)
```
SSH to vm, kill the stress process and uninstall stress-ng 

@param vmId Id of the VM you want to kill the stress process on