// DeploymentService acts as a layer to contain the business logic of the autoscaler, including the delta scaling
// calculation, and handling communication with the api layer.
package main

import (
	"log"
	"math"
)

// A DeploymentConfig contains the data for a deployment's configuration.
type DeploymentConfig struct {
	ServersPerHost int     `json:"servers_per_host"`
	TargetFreePct  float64 `json:"target_free_pct"`
	Id             string  `json:"id"`
}

// A Status contains the data for the current status of an active deployment.
type Status struct {
	CurrentHosts int `json:"current_hosts"`
	TotalServers int `json:"total_servers"`
	FreeServers  int `json:"free_servers"`
}

// CheckStatusAndScale checks the status of the given deployment and scales it, if necessary.
func CheckStatusAndScale(config DeploymentConfig) {
	status, err := GetDeploymentStatus(config.Id)
	if err != nil {
		log.Printf("CheckStatusAndScale: error getting deployment status, skipping deployment: %v", err)
		return
	}

	delta := calculateDelta(config, status)

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

	if config.ServersPerHost < 1 || config.TargetFreePct < 0 || status.FreeServers < 0 || status.TotalServers < 0 {
		log.Printf("calculateDelta: cannot calculate delta for invalid deployment, skipping: config: { %v }, status: { %v }", config, status)
		return delta
	}

	targetFreePct := float64(config.TargetFreePct) / 100.0
	targetFreePct = math.Round(targetFreePct*100) / 100
	freePct := float64(status.FreeServers) / float64(status.TotalServers)

	if freePct != targetFreePct {
		targetBusyPct := 1.0 - targetFreePct
		targetBusyPct = math.Round(targetBusyPct*100) / 100
		busyServerCount := status.TotalServers - status.FreeServers

		// The following uses the inferred target busy percentage and current count of busy servers
		// to calculate the new targetServerCount. i.e. busyServerCount is targetBusyPct of the targetServerCount.
		// We need the target server count in order to calculate the target number of free servers.
		targetServerCount := int(float64(busyServerCount) / targetBusyPct)
		targetFreeServerCount := int(math.Ceil(float64(targetServerCount) * targetFreePct))
		diff := targetFreeServerCount - status.FreeServers

		// Because we're scaling the deployment by hosts, which each contain the configured number of servers, we need to
		// account for that in setting the delta
		delta = int(math.Ceil(float64(diff) / float64(config.ServersPerHost)))
	}

	return delta
}
