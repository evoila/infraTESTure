package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	DeploymentName string `yaml:"deployment_name"`
	TestRepo string `yaml:"test_repo"`
	Service Service `yaml:"service"`
	Testing Testing `yaml:"testing"`
	Bosh Bosh `yaml:"bosh"`
}

type Service struct {
	Name string `yaml:"name"`
	Port int `yaml:"port"`
	Credentials Credentials `yaml:"credentials"`
}

type Credentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Certificate string `yaml:"certificate"`
	Token string `yaml:"token"`
}

type Testing struct {
	Infrastructure string `yaml:"infrastructure"`
	Tests []Test `yaml:"tests"`
}

type Test struct {
	Name string `yaml:"name"`
}

type Bosh struct {
	UaaUrl string `yaml:"uaa_url"`
	DirectorUrl string `yaml:"director_url"`
	UaaClient string `yaml:"uaa_client"`
	UaaClientSecret string `yaml:"uaa_client_secret"`
	Ca string `yaml:"ca"`
}

func LoadConfig(path string) (*Config, error) {
	config := Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Fatalf("[ERROR]: %v", err)
		return nil, err
	}

	return &config, nil
}