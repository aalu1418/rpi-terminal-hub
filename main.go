package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strings"

	"github.com/aalu1418/rpi-terminal-hub/rpio/infrared"
	"github.com/aalu1418/rpi-terminal-hub/services"
	"github.com/aalu1418/rpi-terminal-hub/services/alerts"
	"github.com/aalu1418/rpi-terminal-hub/services/clock"
	"github.com/aalu1418/rpi-terminal-hub/services/connectivity"
	"github.com/aalu1418/rpi-terminal-hub/services/metrics"
	"github.com/aalu1418/rpi-terminal-hub/services/server"
	"github.com/aalu1418/rpi-terminal-hub/services/vacuum"
	"github.com/aalu1418/rpi-terminal-hub/services/weather"
	"github.com/aalu1418/rpi-terminal-hub/types"
	log "github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
)

var (
	OWMKey     string
	IRRecorder bool
	IREmitTest bool
	ClockDemo  bool
)

func init() {
	flag.StringVar(&OWMKey, "owm", os.Getenv(types.OWM_ENVVAR), "pass in openweathermap api key")
	flag.BoolVar(&IRRecorder, "record-ir", false, "run IR recorder")
	flag.BoolVar(&IREmitTest, "emit-ir", false, "run IR emit test")
	flag.BoolVar(&ClockDemo, "clock-demo", false, "run clock demo")
	flag.Parse()
}

func main() {
	if IRRecorder {
		if err := infrared.NewRecorder(rpio.Pin(types.IR_RECEIVER)); err != nil {
			log.Fatalf("[IR] recorder error: %s", err)
		}
		return
	}

	if IREmitTest {
		if err := infrared.NewEmitter(rpio.Pin(types.IR_EMITTER), vacuum.COMMAND_STARTSTOP); err != nil {
			log.Fatalf("[IR] emit test error: %s", err)
		}
		return
	}

	if ClockDemo {
		if err := clock.Demo(rpio.Pin(types.CLOCK_PIN)); err != nil {
			log.Fatalf("[CLOCK] demo error: %s", err)
		}
		return
	}

	if OWMKey == "" {
		log.Fatalf("missing openweathermap api key - pass via %s env or --owm flag", types.OWM_ENVVAR)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	system := make(chan os.Signal, 1)
	signal.Notify(system, os.Interrupt)

	// start message queue
	messages := types.NewQueue()
	defer close(messages)

	// start up services
	clients := []types.Service{
		server.New(messages),
		metrics.New(messages),
		connectivity.New(messages),
		weather.New(messages, OWMKey),
		alerts.NewNWS(messages),
		vacuum.New(messages),
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
