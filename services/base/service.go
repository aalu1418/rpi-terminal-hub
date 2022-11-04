package base

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aalu1418/rpi-terminal-hub/types"
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
		name:         name,
		incomingMsg:  make(chan types.Message),
		outgoingMsg:  make(chan types.Message),
		stop:         make(chan struct{}),
		tickInterval: interval,
		onTick:       onTick,
		processMsg:   processMsg,
	}
}

func (s service) Name() string {
	return s.name
}

func (s *service) Start(_ context.Context) error {
	if s.started {
		return fmt.Errorf("service already started")
	}

	go s.run()
	s.wg.Add(1)

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
			s.processMsg(m)
		case <-tick.C:
			s.outgoingMsg <- s.onTick()
		}
	}
}