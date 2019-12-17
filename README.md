# infraTESTure
Infra Tests? Sure!

## Table of Content

1. [Introduction](#introduction)
2. [Prerequisites](#prerequisites)
3. [Usage](#usage)
4. [Bring Your Own Tests](#bring-your-own-tests)
5. [Troubleshooting](#troubleshooting)
6. [Planned Features](#planned-features)

---

### Introduction

This project is a simple framework written in Go that enables users to easily test deployments, services & infrastructures with either predefined or self-written tests, which can be exchanged in notime. It currently works with [BOSH](https://bosh.io/docs/) deployments only, but is already in work for other infrastructures.

### Prerequisites

The following tools are necessary in order to use this testing framework correctly. The version numbers mentioned in the brackets are not mandatory but simply the versions i used when developing. Other, especially older versions may cause problems.

* [Go](https://golang.org/) (version go1.13.1)
* [Git](https://git-scm.com/) (version 2.23.0)

### Usage

After running `go install` you are ready to go. You can run `infraTESTure --version` in order to check if everything was set up correctly. You should see 

```
infraTESTure CLI version 0.0.1
```

After that you can run `infraTESTure -h` in order to check all possible commands.

```
COMMANDS:
   info, i  Information about what tests are enabled for what services
   run, r   Run tests based on a given configuration file
   help, h  Shows a list of commands or help for one command

```

The info as well as the run command come with one ore more flags you **can** or **must** set.

| Command   | Flags | Description | Flag required |
| --------- | ----------------------------- | ----------------------------------------------------------------------------------------------------- | ------------- |
| info, i   | --repository, -r <br> --tag, -t | URL to the Repository from which you want to get test information <br> Specific Tag from which commit you want to get the test information | yes <br> no | 
| run, r    | --config, -c  <br> --edit, -e <br> --override, -o | Path to the configuration file <br> Tells the tool if you want to edit the test code before running it <br> Overrides existing repository on your computer| yes <br> no <br> no  |

⚠️ **Note:** Only repositories with a specific schemed go code can be used with the `infraTESTure info` command. For more information see [Bring your own tests](#bring-your-own-tests)

If you run the info command with our test repository with predefined redis tests like `infraTESTure info -r https://github.com/evoila/infra-tests`the output should like similar to:

```
2019/12/10 08:30:33 ├── redis
2019/12/10 08:30:33 │   ├── CPU Load
2019/12/10 08:30:33 │   ├── Failover
2019/12/10 08:30:33 │   ├── Health
2019/12/10 08:30:33 │   ├── Info
2019/12/10 08:30:33 │   ├── Network Delay
2019/12/10 08:30:33 │   ├── Package Loss
2019/12/10 08:30:33 │   ├── RAM Load
2019/12/10 08:30:33 │   ├── Service
2019/12/10 08:30:33 │   ├── Storage
```

Now that you have the information about which tests are available for which services you are ready to create the `configuration.yml`:

```yaml
deployment_name: my-test-deployment
github:
  test_repo: https://github.com/evoila/infra-tests
  tag: v1.0
  saving_location: /Users/me/Desktop
  repo_name: my-tests
service:
  name: redis
  port: 6379
  credentials:
    username: someUsername
    password: somePassword
    certificate: someCertificate
    token: someToken
testing:
  infrastructure: bosh
  tests:
  - name: info
  - name: health
  - name: service
  - name: failover
  - name: storage
    properties:
      path: /path/to/persistent/disk
  - name: package loss
    properties:
      directorIp: 127.0.0.1
  - name: network delay
    properties:
      directorIp: 127.0.0.1
  - name: cpu load
  - name: ram load
bosh:
  uaa_url: https://127.0.0.1:8443
  director_url: https://127.0.0.1:25555
  uaa_client: admin
  uaa_client_secret: adminPassword
```

⚠️ **Note:** The test execution order is equals to the order the tests are defined in the `configuration.yml`. So the info test is executed first, health second, service third etc.

| Field                                  | Description   | Required |
| -------------------------------------- | ----------------- | ----- |
| deployment_name | Name of the deployment | Yes |
| github.test_repo| URL of the github repository containing the tests| Yes |
| github.tag | Specific Tag you want to clone the repository from | No | 
| github.saving_location | Path on your computer where the github repo is going to be saved to | Yes |
| github.repo_name | Describes under which directory name the repository is saved on your computer | Yes|
| service<span>.name| Name of the service you want to test | Yes |
| service.port | Port of the service you want to test | Yes 
| service.credentials.username | Usernamen for the service | Depends on service |
| service.credentials.password | Password for the service | Depends on service | 
| service.credentials.certificate | Certificate for the service | Depends on service |
| service.credentials.token | Token for the service | Depends on service |
| testing.infrastructure | Name of the infrastructure your services are running on | Yes |
| testing.tests | List of tests you want to run | No |
| bosh.uaa_url | UAA URL of the bosh deployment | Yes |
| bosh.director_url | Director URL of the bosh deployment | Yes |
| bosh.uaa_client | Usernamen of the bosh UAA client | Yes |
| bosh.uaa_client_secret | Password for the bosh UAA client | Yes|

As you can see some tests also have a field `properties`. It can be important to individually configure some tests. For example we need the bosh director ip for the traffic shaping tests, to exclude it from the traffic shaping rules. Imagine simulating 100% package loss on a VM. The bosh director would think the VM is unresponsive and would try to repair it.
The good thing with these `properties`is that you can add whatever property you want. In your test simply run

```go
getTestProperties(config, "<test_name>")["<property_name>"]
```

After setting up the configuration.yml you should now be able to run your tests with `infraTESTure run -c /path/to/configuration.yml`. The output should look similar to this:

```
2019/10/14 15:12:10 [INFO] Cloning repository from https://github.com/evoila/infra-tests
2019/10/14 15:12:11 [INFO] Building go plugin from directory /Users/me/Desktop/infra-tests/redis
2019/10/14 15:12:12 [INFO] Loading go plugin...

##### Health Test #####
2019/10/14 15:12:17 [INFO] Checking process state for every VM of Deployment my-test-deployment...
...
...
...
```

### Bring Your Own Tests

This project was designed as a full community driven and generic testing framework for infrastructures and services, which means that you are able to use your very own tests. When writing the code, there are some restrictions you have to follow in order to make the framework work with your tests.

#### Project Structure

When creating the project the first important but easy to handle restriction is the project structure. For every service you want to implement tests for you have to create a folder in the root directory of the project. In this folder you can create several go files with your code.

⚠️ **Note:** The package information of **every** go file has to be `package main` in order to make the go plugin work.

Lets say you want to create tests for MongoDB and Redis, your project structure has to look like this:

```
├──infra-tests (root)
│   ├── mongodb
│   │ 	├── firstMongodbFile.go
│   │ 	├── secondMongodbFile.go
│   │ 	├── thirdMongodbFile.go
│   ├── redis
│   │ 	├── firstRedisFile.go
│   │ 	├── secondRedisFile.go
│   │ 	├── thirdRedisFile.go
```

Remember that you have to adjust your configuration.yml, where `github.test_repo` is now the URL to your own repository and `service.name` is equal to one of the directories you created (mongodb, redis...).

#### Function Signatures

First of all run `go get github.com/evoila/infraTESTure`

Every function you want to be executed on runtime must have a specific signature in order to be found by the go plugin package.

```go
func MyTest(config *config.Config, infrastructure infrastructure.Infrastructure) bool { ... } 
```

with a boolean return value telling you whether the test was successful (true) or not (false) and parameters `Config` and `Infrastructure` imported from

```go
"github.com/evoila/infraTESTure/config"
"github.com/evoila/infraTESTure/infrastructure"
```

#### Annotations

Last but not least you have to annotate the function correctly. Since there are not official annotations in Go we solved this problem by simply using comments. This comment has to be exactly above the function and must look like

```go
// @Test
func MyTest(...) bool { ... }
```

You could now add 

```
testing:
  tests:
  - name: test
```

to your `configuration.yml` (where tests<span>.name is equal to the annotation) in order to execute this test when running `infraTESTure run`. Combining all these restrictions results in a file that should look like [this](https://github.com/evoila/infra-tests/blob/master/test/test.go).

### Troubleshooting

coming soon...

### Planned Features

coming soon...