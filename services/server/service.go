package server

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/aalu1418/rpi-terminal-hub/services/base"
	"github.com/aalu1418/rpi-terminal-hub/services/vacuum"
	"github.com/aalu1418/rpi-terminal-hub/services/weather"
	"github.com/aalu1418/rpi-terminal-hub/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

//go:embed static index.html
var static embed.FS

type service struct {
	types.Service
	server   *http.Server
	template *template.Template
	data     outputData
	lock     sync.RWMutex
	out      chan<- types.Message
}

type outputData struct {
	Weather types.WeatherParsed
	Alerts  string
	Vacuum  string
}

func New(outgoingMsg chan<- types.Message) types.Service {
	var s service
	s.Service = base.New(outgoingMsg, types.WEBSERVER, types.INFINITE_TIME, s.onTick, s.processMsg)
	s.server = &http.Server{Addr: types.WEBSERVER_ADDRESS}
	s.template = template.Must(template.ParseFS(static, "index.html"))
	s.data.Weather = weather.EMPTY
	s.out = outgoingMsg
	return &s
}

// custom Start
func (s *service) Start(ctx context.Context) error {
	if err := s.Service.Start(ctx); err != nil {
		return err
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("pong")); err != nil {
			log.Errorf("server.ping: %s", err)
		}
	})

	http.Handle("/metrics", promhttp.Handler())

	http.Handle("/static/", http.FileServer(http.FS(static)))

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		s.lock.RLock()
		defer s.lock.RUnlock()
		if err := s.template.Execute(w, s.data); err != nil {
			log.Errorf("server.weatherL %s", err)
		}
	})

	http.HandleFunc("/vacuum", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			if _, err := w.Write([]byte("invalid method")); err != nil {
				log.Errorf("server.vacuum.method %s", err)
			}
			return
		}

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				log.Errorf("server.vacuum.ioReadAll %s", err)
			}
			return
		}

		cmd := string(bytes)
		var exists bool
		for _, v := range vacuum.Calls {
			exists = exists || v == cmd
		}
		if !exists {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte(fmt.Sprintf("invalid command: %s", cmd))); err != nil {
				log.Errorf("server.vacuum.exists %s", err)
			}
			return
		}

		s.out <- types.Message{
			From: s.Name(),
			To:   types.VACUUM,
			Data: cmd,
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(fmt.Sprintf("%s cmd emitted", strings.ToUpper(cmd)))); err != nil {
			log.Errorf("server.vacuum.output %s", err)
		}
	})

	// start server in go routine
	go func() {
		// ErrServerClosed on graceful close
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("[CRITICAL] web-server fatally errored: %s", err)
		}
	}()

	return nil
}

// custom Stop
func (s *service) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	return s.Service.Stop(ctx)
}

func (s *service) processMsg(m types.Message) {
	if m.To != types.WEBSERVER {
		log.Errorf("[SERVER] recieved incorrect message: %+v", m)
		return
	}

	switch {
	case m.From == types.WEATHER:
		weatherData, ok := m.Data.(types.WeatherParsed)
		if !ok {
			log.Errorf("[SERVER] could not parse weather message: %+v", m)
			return
		}
		s.lock.Lock()
		defer s.lock.Unlock()
		s.data.Weather = weatherData
		return
	case m.From == types.NWS:
		alert, ok := m.Data.(string)
		if !ok {
			log.Errorf("[SERVER] could not parse NWS alert message: %+v", m)
			return
		}
		s.lock.Lock()
		defer s.lock.Unlock()
		s.data.Alerts = alert
		return
	case m.From == types.VACUUM:
		status, ok := m.Data.(string)
		if !ok {
			log.Errorf("[SERVER] could not parse vacuum status message: %+v", m)
			return
		}
		s.lock.Lock()
		defer s.lock.Unlock()
		s.data.Vacuum = status
		return
	default:
		log.Infof("[SERVER] received: %+v", m)
	}
}

// called once at the very beginning (and after INFINITE_TIME)
func (s *service) onTick() types.Message {
	return types.Message{
		To:   types.POSTOFFICE,
		Data: fmt.Sprintf("[ALIVE] %s", s.Name()),
	}
}
