package deployments

import (
	"encoding/json"
	"fmt"
)

type discoveredDeployments struct {
	Configs []Config `json:"deployments"`
}

func Discover() []Config {
	fmt.Println("Discovering Deployments")
	//get /discover
	body := DiscoverDeployments()

	var discoveredDeploys discoveredDeployments
	err := json.Unmarshal(body, &discoveredDeploys)
	if err != nil {
		fmt.Printf("error: could not unmarshal: %s\n", err)
	}

	return discoveredDeploys.Configs
}

func NeedsScaling(config Config) bool {
	status := status(config.Id)
	if status.CurrentHosts > 0 && status.TotalServers > 0 && status.FreeServers > 0 {
		return true
	}
	return false
}

func status(id string) Status {
	fmt.Printf("Getting Status for deployment: %v\n", id)
	body := DeploymentStatus(id)

	var status Status
	err := json.Unmarshal(body, &status)
	if err != nil {
		fmt.Printf("error: could not unmarshal: %s\n", err)
	}

	return status
}

func Scale(config Config) {
	//post /scale with delta and deployment_id
	ScaleDeployment(5, config.Id)
	fmt.Printf("Scaling: %v\n", config)
}
