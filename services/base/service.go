package base

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aalu1418/rpi-terminal-hub/types"
	log "github.com/sirupsen/logrus"
)

type service struct {
	name         string
	started      bool
	incomingMsg  chan types.Message
	outgoingMsg  chan types.Message
	wg           sync.WaitGroup
	stop         chan struct{}
	tickInterval time.Duration
	onTick       types.OnTick
	processMsg   types.ProcessMsg
}

func New(name string, interval time.Duration, onTick types.OnTick, processMsg types.ProcessMsg) types.Service {
	return &service{
		name:         strings.ToLower(name),
		started:      false,
		incomingMsg:  make(chan types.Message),
		outgoingMsg:  make(chan types.Message),
		wg:           sync.WaitGroup{},
		stop:         make(chan struct{}),
		tickInterval: interval,
		onTick:       onTick,
		processMsg:   processMsg,
	}
}

func (s *service) Name() string {
	return s.name
}

func (s *service) Start(_ context.Context) error {
	if s.started {
		return fmt.Errorf("service already started")
	}

	go s.run()
	s.wg.Add(1)
	s.started = true
	return nil
}

func (s *service) Healthy() bool {
	return s.started
}

func (s *service) ExtRead() <-chan types.Message {
	return s.outgoingMsg
}

func (s *service) ExtWrite() chan<- types.Message {
	return s.incomingMsg
}

func (s *service) Stop() error {
	if !s.started {
		return fmt.Errorf("service not started")
	}

	close(s.stop)

	s.wg.Wait()
	close(s.incomingMsg)
	close(s.outgoingMsg)
	s.started = false
	return nil
}

func (s *service) run() {
	tick := time.NewTicker(s.tickInterval)
	defer tick.Stop()
	defer s.wg.Done()

	for {
		select {
		case <-s.stop:
			return
		case m := <-s.incomingMsg:
			// discard message if To does not match name or From is not included
			if m.To != s.name || m.From == "" {
				log.Warnf("service (%s) received invalid msg: %+v", s.name, m)
				continue // continue loop
			}
			s.processMsg(m)
		case <-tick.C:
			msg := s.onTick()
			msg.From = s.name // always use name as from
			s.outgoingMsg <- msg
		}
	}
}

// XXXBaseImplementations for testing only
// provides the needed interfaces to construct
func XXXNewBaseImplementation(t *testing.T, name string, processMsg types.ProcessMsg) types.Service {
	duration := time.Second
	onTick := func() types.Message {
		return types.Message{
			To:   "placeholder",
			Data: time.Now(),
		}
	}
	return New(name, duration, onTick, processMsg)
}