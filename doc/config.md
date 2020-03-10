# config
--
    import "github.com/evoila/infraTESTure/config"


## Usage

#### type Bosh

```go
type Bosh struct {
	UaaUrl          string `yaml:"uaa_url"`
	DirectorUrl     string `yaml:"director_url"`
	UaaClient       string `yaml:"uaa_client"`
	UaaClientSecret string `yaml:"uaa_client_secret"`
	Ca              string `yaml:"ca"`
}
```


#### type Config

```go
type Config struct {
	DeploymentName string  `yaml:"deployment_name"`
	Github         Github  `yaml:"github"`
	Service        Service `yaml:"service"`
	Testing        Testing `yaml:"testing"`
	Bosh           Bosh    `yaml:"bosh"`
}
```


#### func  LoadConfig

```go
func LoadConfig(path string) (*Config, error)
```
Parse the content of a configuration file to the above go structs 
@param path Path to the configuration file 
@return config Config struct containing the information from the configuration file

#### type Credentials

```go
type Credentials struct {
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Certificate string `yaml:"certificate"`
	Token       string `yaml:"token"`
}
```


#### type Github

```go
type Github struct {
	TestRepo       string `yaml:"test_repo"`
	Tag            string `yaml:"tag"`
	SavingLocation string `yaml:"saving_location""`
	RepoName       string `yaml:"repo_name""`
}
```


#### type Service

```go
type Service struct {
	Name        string      `yaml:"name"`
	Port        int         `yaml:"port"`
	Credentials Credentials `yaml:"credentials"`
}
```


#### type Test

```go
type Test struct {
	Name       string            `yaml:"name"`
	Properties map[string]string `yaml:"properties"`
}
```


#### type Testing

```go
type Testing struct {
	Infrastructure string `yaml:"infrastructure"`
	Tests          []Test `yaml:"tests"`
}
```
