// Autoscaler is a service that automates the retrieval of a list of deployments on a server as well as the scaling
// of hosts on server deployments to maintain the configured target percentage of free servers on the deployment.
package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// The number of seconds to wait between each check.
const (
	secondsBetweenDiscovery     = 60
	secondsBetweenCheckAndScale = 30
)

// Main operates as the entry point to the program.
func main() {
	// Set up logging.
	file, err := os.OpenFile("autoscaler.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	// Set up graceful exit.
	quitSignalChannel := make(chan os.Signal, 1)
	signal.Notify(quitSignalChannel, syscall.SIGINT, syscall.SIGTERM)
	shutdownChannel := make(chan bool)
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	autoscale(shutdownChannel, waitGroup)

	// Wait for, and then handle graceful exit.
	<-quitSignalChannel
	log.Println("Received quit. Sending shutdown then waiting on goroutines...")
	shutdownChannel <- true
	waitGroup.Wait()
	log.Println("Exiting.")
}

// Autoscale runs the goroutines that discover deployments and check status/scale each deployment.
func autoscale(shutdownChannel chan bool, waitGroup *sync.WaitGroup) {
	log.Println("Starting autoscaler...")
	activeDeployments := make(chan []DeploymentConfig)

	// Discovery goroutine to handle the discover deployments interval
	go func(shutdownChannel chan bool, waitGroup *sync.WaitGroup) {
		defer waitGroup.Done()
		defer close(activeDeployments)

		ticker := time.NewTicker(secondsBetweenDiscovery * time.Second)
		defer ticker.Stop()

		for ; true; <-ticker.C {
			select {
			case <-shutdownChannel:
				return
			default:
			}

			discoveredDeployments, err := GetDeployments()
			if err != nil {
				log.Printf("Autoscale: error getting deployments, skipping interval: %v", err)
				continue
			} else if len(discoveredDeployments) == 0 {
				log.Printf("Autoscale: received empty list of deployments, skipping interval: %v", err)
				continue
			}

			activeDeployments <- discoveredDeployments
			timer := time.After(secondsBetweenCheckAndScale * time.Second)
			<-timer
			activeDeployments <- discoveredDeployments // Sends to the activeDeployments channel twice as the status check and scaling occurs every 30 sec
		}
	}(shutdownChannel, waitGroup)

	// Status/Scale goroutine to handle the deployment status and scale interval
	go func(shutdownChannel chan bool, waitGroup *sync.WaitGroup) {
		defer waitGroup.Done()

		for {
			select {
			case <-shutdownChannel:
				return
			default:
			}

			activeDeploys := <-activeDeployments
			for _, deployment := range activeDeploys {
				go CheckStatusAndScale(deployment)
			}
		}
	}(shutdownChannel, waitGroup)
}
