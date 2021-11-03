package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/automaxprocs/maxprocs"
)

var build = "develop"

func main() {

	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotas
	if _, err := maxprocs.Set(); err != nil {
		fmt.Println("Maxprocs: %w", err)
		os.Exit(1)
	}
	cores := runtime.GOMAXPROCS(0)

	log.Printf("Starting service...build[%s] CPU[%d]", build, cores)
	defer log.Println("Service ended.")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("Stopping service...", build)
}
