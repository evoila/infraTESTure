package bosh

import (
	"github.com/briandowns/spinner"
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	"github.com/evoila/infraTESTure/config"
	"log"
	"time"
)

func IsDeploymentRunning(config *config.Config) bool {
	spin := spinner.New(spinner.CharSets[33], 100*time.Millisecond)

	log.Printf("[INFO] Checking process state for every VM of Deployment %v...", config.DeploymentName)

	spin.Start()

	director, err := buildDirector(config)
	if err != nil {
		log.Fatal(err)
	}

	// See director/interfaces.go for a full list of methods.
	deps, err := director.FindDeployment(config.DeploymentName)
	if err != nil {
		log.Fatal(err)
	}

	var vmVitals []boshdir.VMInfo
	vmVitals, err = deps.VMInfos()

	if err != nil {
		log.Fatal(err)
	}

	healthy := true

	spin.Stop()

	for _, vmVital := range vmVitals {
		log.Printf("[INFO] %v - %v/%v - %v", *vmVital.Index, vmVital.JobName, vmVital.ID, vmVital.ProcessState)

		if vmVital.ProcessState != "running" {
			healthy = false
		}
	}

	return healthy
}