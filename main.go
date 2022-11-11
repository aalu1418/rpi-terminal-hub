package main

import (
	"context"
	"os"
	"os/signal"
	"strings"

	"github.com/aalu1418/rpi-terminal-hub/services"
	"github.com/aalu1418/rpi-terminal-hub/services/connectivity"
	"github.com/aalu1418/rpi-terminal-hub/services/metrics"
	"github.com/aalu1418/rpi-terminal-hub/services/server"
	"github.com/aalu1418/rpi-terminal-hub/services/weather"
	"github.com/aalu1418/rpi-terminal-hub/types"
	log "github.com/sirupsen/logrus"
)

var OWMKey string

func init() {
	OWMKey = os.Getenv("OWM_KEY")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	system := make(chan os.Signal, 1)
	signal.Notify(system, os.Interrupt)

	// start message queue
	messages := make(chan types.Message)
	defer close(messages)

	// start up services
	clients := []types.Service{
		server.New(messages),
		metrics.New(messages),
		connectivity.New(messages),
		weather.New(messages, OWMKey),
	}

	// start up post office for message sorting
	postOffice := services.NewPostOffice(messages, clients)

	// start postOffice first
	startStop := []types.BaseService{postOffice}
	for i := range clients {
		startStop = append(startStop, clients[i])
	}

	// start up all services
	for i := range startStop {
		startCtx, cancel := context.WithTimeout(ctx, types.DEFAULT_TIMEOUT)
		if err := startStop[i].Start(startCtx); err != nil {
			log.Panicf("service (%s) failed to start: %s", startStop[i].Name(), err)
		}
		cancel()
	}

	var names []string
	for i := range clients {
		names = append(names, clients[i].Name())
	}
	log.Infof("all services started: %s", strings.Join(names, ","))

	// wait for exit system to interrupt
	<-system
	log.Info("shutting down services")

	// stop all clients
	// stop post office first to avoid sending on closed channel
	for i := range startStop {
		stopCtx, cancel := context.WithTimeout(ctx, types.DEFAULT_TIMEOUT)
		if err := startStop[i].Stop(stopCtx); err != nil {
			log.Warnf("service (%s) failed to stop properly: %s", startStop[i].Name(), err)
		}
		cancel()
	}

	log.Info("all services shut down")
}
