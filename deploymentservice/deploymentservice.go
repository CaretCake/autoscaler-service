package deploymentservice

import "fmt"

type deploymentConfig struct {
	serversPerHost int
	targetFreePct  int
}

func DiscoverDeployments() {
	fmt.Println("Discovering Deployments")
}

func Status() {
	fmt.Println("Getting Status for: ")
}

func Scale() {
	fmt.Println("Scaling: ")
}
