package main

import (
	"fmt"
	"log"
	"math"
)

// A DeploymentConfig contains the data for a deployment's configuration.
type DeploymentConfig struct {
	ServersPerHost int     `json:"servers_per_host"`
	TargetFreePct  float32 `json:"target_free_pct"`
	Id             string  `json:"id"`
}

// A Status contains the data for the current status of an active deployment.
type Status struct {
	CurrentHosts int `json:"current_hosts"`
	TotalServers int `json:"total_servers"`
	FreeServers  int `json:"free_servers"`
}

// Discover returns a list of DeploymentConfigs for all active deployments.
func Discover() ([]DeploymentConfig, error) {
	discoveredDeploys, err := GetDeployments()
	if err != nil {
		return nil, fmt.Errorf("Discover: %v", err)
	}

	return discoveredDeploys, nil
}

// CheckStatusAndScale checks the status of the given deployment and scales it, if necessary.
func CheckStatusAndScale(config DeploymentConfig) {
	status, err := GetDeploymentStatus(config.Id)
	if err != nil {
		log.Printf("CheckStatusAndScale: error getting deployment status, skipping deployment: %v", err)
		return
	}

	delta := 0
	delta = calculateDelta(config, status)

	if delta != 0 {
		err = ScaleDeployment(config.Id, delta)
		if err != nil {
			log.Printf("CheckStatusAndScale: scaling of deployment { %v } failed, skipping deployment : %v", config, err)
		}
	}
}

// CalculateDelta returns the delta by which to scale the number of hosts on the given deployment to maintain
// the target percentage of free servers. Note that the returned int value may be positive or negative.
func calculateDelta(config DeploymentConfig, status Status) int {
	delta := 0
	targetFreePct := float64(config.TargetFreePct) / 100.0
	targetFreePct = math.Round(targetFreePct*100) / 100
	freePercent := float64(status.FreeServers) / float64(status.TotalServers)

	if freePercent != targetFreePct {
		targetBusyPct := 1.0 - targetFreePct
		targetBusyPct = math.Round(targetBusyPct*100) / 100
		busyServerCount := status.TotalServers - status.FreeServers

		// The following uses the inferred target busy percentage and current number of busy servers
		// to calculate the new targetServerCount. i.e. busyServerCount is targetBusyPct of the targetServerCount
		targetServerCount := int(float64(busyServerCount) / targetBusyPct)
		targetFreeServerCount := int(math.Ceil(float64(targetServerCount) * targetFreePct))
		diff := targetFreeServerCount - status.FreeServers

		delta = int(math.Ceil(float64(diff) / float64(config.ServersPerHost)))
	}

	return delta
}
