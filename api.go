// Api contains all of the calls out to the deployments api and handles reading, deserialization / serialization,
// as well as errors from the api.
package main

import (
	"fmt"
)

const apiURL = "http://127.0.0.1:5000"

// GetDeployments makes a GET request to return a list of DeploymentConfigs representing all active deployments.
func GetDeployments() ([]DeploymentConfig, error) {
	body, err := TryGet(apiURL + "/discover")
	if err != nil {
		return nil, fmt.Errorf("GetDeployments: %s", err)
	}

	type discoveredDeployments struct {
		Configs []DeploymentConfig `json:"deployments"`
	}
	var discoveredDeploys discoveredDeployments
	err = TryUnmarshalJSON(body, &discoveredDeploys)
	if err != nil {
		return nil, fmt.Errorf("GetDeployments: %s", err)
	}

	return discoveredDeploys.Configs, nil
}

// GetDeploymentStatus makes a GET request to return the Status of the given deployment.
func GetDeploymentStatus(id string) (Status, error) {
	body, err := TryGet(apiURL + "/status/" + id)
	if err != nil {
		return Status{}, fmt.Errorf("GetDeploymentStatus: %s", err)
	}

	var s Status
	err = TryUnmarshalJSON(body, &s)
	if err != nil {
		return Status{}, fmt.Errorf("GetDeploymentStatus: %w", err)
	}

	return s, nil
}

// ScaleDeployment makes a POST request to scale the number of hosts by delta on the given deployment.
func ScaleDeployment(id string, delta int) error {
	type ScalePayload struct {
		DeploymentId string `json:"deployment_id"`
		Delta        int    `json:"delta"`
	}
	payload := ScalePayload{
		DeploymentId: id,
		Delta:        delta,
	}
	_, err := TryPostJSON(apiURL+"/scale", payload)
	if err != nil {
		return fmt.Errorf("ScaleDeployment: %s", err)
	}

	return nil
}
