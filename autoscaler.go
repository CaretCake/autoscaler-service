package main

import (
	"autoscaler/deploymentservice"
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
		done <- true
	}()

	go func() {
		for {
			deploymentservice.DiscoverDeployments()
			time.Sleep(time.Minute)
		}
	}()

	go func() {
		for {
			deploymentservice.Status()
			deploymentservice.Scale()
			time.Sleep(time.Second * 30)
		}
	}()

	<-done
	fmt.Println("Exiting")
}
