package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/aalu1418/rpi-terminal-hub/services"
	"github.com/aalu1418/rpi-terminal-hub/types"
	log "github.com/sirupsen/logrus"
)

const (
	TIMEOUT = 5 * time.Second
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	system := make(chan os.Signal, 1)
	signal.Notify(system, os.Interrupt)

	// start message queue
	messages := make(chan types.Message)
	defer close(messages)

	// start up services
	clients := []types.Service{}

	// start up post office for message sorting
	postOffice := services.NewPostOffice(messages, clients)

	// start postOffice first
	startStop := []types.BaseService{postOffice}
	for i := range clients {
		startStop = append(startStop, clients[i])
	}

	// start up all services
	for i := range startStop {
		startCtx, cancel := context.WithTimeout(ctx, TIMEOUT)
		if err := startStop[i].Start(startCtx); err != nil {
			log.Fatalf("service (%s) failed to start: %s", startStop[i].Name(), err)
		}
		cancel()
	}

	// wait for exit system to interrupt
	<-system
	log.Info("shutting down services")

	// stop all clients
	// stop post office first to avoid sending on closed channel
	for i := range startStop {
		if err := startStop[i].Stop(); err != nil {
			log.Warnf("service (%s) failed to stop properly: %s", startStop[i].Name(), err)
		}
	}

	log.Info("all services shut down")
}
