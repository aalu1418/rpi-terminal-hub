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
	incomingMsg  chan types.Message   // provided to messaging service to send data back
	outgoingMsg  chan<- types.Message // provided externally by messaging service
	wg           sync.WaitGroup
	stop         chan struct{}
	tickInterval time.Duration
	onTick       types.OnTick
	processMsg   types.ProcessMsg
}

func New(outgoingMsg chan<- types.Message, name string, interval time.Duration, onTick types.OnTick, processMsg types.ProcessMsg) types.Service {
	return &service{
		name:         strings.ToLower(name),
		started:      false,
		incomingMsg:  types.NewQueue(),
		outgoingMsg:  outgoingMsg,
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

	go s.runOutgoing()
	go s.runIncoming()
	s.wg.Add(2)
	s.started = true
	return nil
}

func (s *service) Healthy() bool {
	return s.started
}

func (s *service) ExtWrite() chan<- types.Message {
	return s.incomingMsg
}

func (s *service) Stop(_ context.Context) error {
	if !s.started {
		return fmt.Errorf("service not started")
	}

	close(s.stop)
	s.wg.Wait()
	close(s.incomingMsg)
	s.started = false
	return nil
}

func (s *service) runOutgoing() {
	tick := time.After(0)
	defer s.wg.Done()

	for {
		select {
		case <-s.stop:
			return
		case <-tick:
			start := time.Now()                                   // track starting time
			msg := s.onTick()                                     // run msg generating function
			msg.From = s.name                                     // always use name as from
			tick = time.After(s.tickInterval - time.Since(start)) // reset tick

			// handle full outoing queue
			if len(s.outgoingMsg) == types.MAX_QUEUE {
				log.Errorf("[CRITICAL] outgoing message queue is full, msg dropped: %+v", msg)
				continue
			}
			s.outgoingMsg <- msg
		}
	}
}

func (s *service) runIncoming() {
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
		}
	}
}

// XXXBaseImplementations for testing only
// provides the needed interfaces to construct
func XXXNewBaseImplementation(t *testing.T, out chan<- types.Message, name string, processMsg types.ProcessMsg) types.Service {
	duration := time.Second
	onTick := func() types.Message {
		return types.Message{
			To:   "placeholder",
			Data: time.Now(),
		}
	}
	return New(out, name, duration, onTick, processMsg)
}