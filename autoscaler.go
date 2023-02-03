package main

import (
	"autoscaler/deployments"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		signal := <-sigs
		fmt.Println(signal)
		// wait group or whatever with the others to ensure that we let them empty their channels?
		done <- true
	}()

	activeDeployments := make(chan []deployments.Config, 1)
	scalableDeployments := make(chan deployments.Config)

	go func() {
		for {
			discoveredDeployments := deployments.Discover()
			activeDeployments <- discoveredDeployments
			activeDeployments <- discoveredDeployments
			time.Sleep(time.Minute)
		}
	}()

	go func() {
		for {
			activeDeploys := <-activeDeployments
			for _, deployment := range activeDeploys {
				if deployments.NeedsScaling(deployment) {
					scalableDeployments <- deployment
				}
			}
			time.Sleep(time.Second * 30)
		}
	}()

	go func() {
		for {
			deployment := <-scalableDeployments
			deployments.Scale(deployment)
		}
	}()

	<-done
	fmt.Println("Exiting")
}
