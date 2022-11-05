package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/aalu1418/rpi-terminal-hub/services"
	"github.com/aalu1418/rpi-terminal-hub/types"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	system := make(chan os.Signal, 1)
	signal.Notify(system, os.Interrupt)

	// start message queue
	messages := make(chan types.Message)

	// start up services
	clients := []types.Service{}

	// start up post office for message sorting
	postOffice := services.NewPostOffice(messages, clients)
	if err := postOffice.Start(ctx); err != nil {
		log.Fatalf("post office failed to start: %s", err)
	}

	// start up all clients

	// wait for exit system to interrupt
	select {
	case <-system:
		// stop all clients

		// stop the post office
	}
	return
}
