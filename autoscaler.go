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
	file, err := os.OpenFile("autoscalerLog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
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
	activeDeployments := make(chan []DeploymentConfig)

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

			discoveredDeployments, err := Discover()
			if err != nil {
				log.Printf("Autoscale: error getting deployments, skipping: %v", err)
				continue
			}
			//log.Println("Discovered deploymentss: ", time.Now())

			activeDeployments <- discoveredDeployments
			timer := time.After(secondsBetweenCheckAndScale * time.Second)
			<-timer
			activeDeployments <- discoveredDeployments // Sends to the activeDeployments channel twice as the status check and scaling occurs every 30 sec
		}
	}(shutdownChannel, waitGroup)

	go func(shutdownChannel chan bool, waitGroup *sync.WaitGroup) {
		defer waitGroup.Done()

		for {
			select {
			case <-shutdownChannel:
				return
			default:
			}

			activeDeploys := <-activeDeployments
			//log.Println("StatusCheck received deployments: ", activeDeploys, " : ", time.Now())
			for _, deployment := range activeDeploys {
				go CheckStatusAndScale(deployment)
			}
		}
	}(shutdownChannel, waitGroup)
}
