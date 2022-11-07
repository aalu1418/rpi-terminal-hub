package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aalu1418/rpi-terminal-hub/services/base"
	"github.com/aalu1418/rpi-terminal-hub/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func init() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("pong")); err != nil {
			log.Errorf("server.ping: %s", err)
		}
	})

	http.Handle("/metrics", promhttp.Handler())
}

type service struct {
	types.Service
	server *http.Server
}

func New(outgoingMsg chan<- types.Message) types.Service {
	var s service
	s.Service = base.New(outgoingMsg, types.WEBSERVER, types.INFINITE_TIME, s.onTick, s.processMsg)
	s.server = &http.Server{Addr: types.WEBSERVER_ADDRESS}
	return &s
}

// custom Start
func (s *service) Start(ctx context.Context) error {
	// start server in go routine
	go func() {
		// ErrServerClosed on graceful close
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("[CRITICAL] web-server fatally errored: %s", err)
		}
	}()

	return s.Service.Start(ctx)
}

// custom Stop
func (s *service) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	return s.Service.Stop(ctx)
}

func (s *service) processMsg(m types.Message) {
	log.Infof("received msg: %+v", m)
}

// called once at the very beginning (and after INFINITE_TIME)
func (s *service) onTick() types.Message {
	return types.Message{
		To:   types.POSTOFFICE,
		Data: fmt.Sprintf("[ALIVE] %s", s.Name()),
	}
}
