package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aalu1418/rpi-terminal-hub/services/base"
	"github.com/aalu1418/rpi-terminal-hub/types"
	log "github.com/sirupsen/logrus"
)

const (
	NAME = "web-server"
)

func init() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
}

type service struct {
	base   types.Service
	server *http.Server
}

func New(outgoingMsg chan<- types.Message) types.Service {
	var s service
	s.base = base.New(outgoingMsg, NAME, types.INFINITE_TIME, s.onTick, s.processMsg)
	s.server = &http.Server{Addr: types.WEBSERVER_ADDRESS}
	return &s
}

func (s *service) Name() string {
	return s.base.Name()
}

func (s *service) Start(ctx context.Context) error {
	// start server in go routine
	go func() {
		// ErrServerClosed on graceful close
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("[CRITICAL] web-server fatally errored: %s", err)
		}
	}()

	return s.base.Start(ctx)
}

func (s *service) Healthy() bool {
	return s.base.Healthy()
}

func (s *service) ExtWrite() chan<- types.Message {
	return s.base.ExtWrite()
}

func (s *service) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	return s.base.Stop(ctx)
}

func (s *service) processMsg(m types.Message) {
	log.Infof("received msg: %+v", m)
}

// called once at the very beginning (and after INFINITE_TIME)
func (s *service) onTick() types.Message {
	return types.Message{
		To:   types.POSTOFFICE,
		Data: fmt.Sprintf("[ALIVE] %s", NAME),
	}
}
