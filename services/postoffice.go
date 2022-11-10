package services

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/aalu1418/rpi-terminal-hub/types"
	log "github.com/sirupsen/logrus"
)

// sorts the messages to the appropriate receiver like the post office
type postOffice struct {
	mailbag   <-chan types.Message
	mailboxes map[string]chan<- types.Message
	started   bool
	stop      chan struct{}
	workers   atomic.Int32
}

func NewPostOffice(incoming <-chan types.Message, clients []types.Service) types.BaseService {
	m := postOffice{
		mailbag:   incoming,
		mailboxes: map[string]chan<- types.Message{},
		stop:      make(chan struct{}),
	}

	// create mailboxes
	for i := range clients {
		address := strings.ToLower(clients[i].Name())

		if _, exists := m.mailboxes[address]; exists {
			log.Panicf("duplicate service names: %s", address)
		}

		m.mailboxes[address] = clients[i].ExtWrite()
	}

	return &m
}

func (m *postOffice) Name() string {
	return types.POSTOFFICE
}

func (m *postOffice) Start(_ context.Context) error {
	if m.started {
		return fmt.Errorf("post office has been started")
	}
	go m.sort()
	m.started = true
	return nil
}

func (m *postOffice) sort() {
	for {
		select {
		case <-m.stop:
			return
		case msg := <-m.mailbag:
			m.workers.Add(1) // try to add worker for new message
			if m.workers.Load() > types.MAX_WORKERS {
				log.Errorf("[CRITICAL] postoffice: max workers reached (%d), consider increasing limits", types.MAX_WORKERS)
			}
			for m.workers.Load() > types.MAX_WORKERS {
				// blocking
				// wait until a worker frees up
			}
			go m.worker(msg)
		}
	}
}

func (m *postOffice) worker(msg types.Message) {
	defer m.workers.Add(-1) // subtract worker when done

	if msg.To == "" || msg.From == "" {
		log.Errorf("post office received invalid message: %+v", msg)
		return
	}

	// handle msg if sent to post office
	if msg.To == types.POSTOFFICE {
		log.Infof("[POSTOFFICE] log: %+v", msg)
		return
	}

	mailbox, exists := m.mailboxes[msg.To]
	if !exists {
		log.Errorf("no valid mailbox for addressee: %+v", msg)
		return
	}

	if len(mailbox) == types.MAX_QUEUE {
		log.Errorf("[CRITICAL] mailbox full, dropping message: %+v", msg)
		return
	}

	mailbox <- msg
}

func (m *postOffice) Stop(_ context.Context) error {
	if !m.started {
		return fmt.Errorf("post office has not been started")
	}
	close(m.stop)
	m.started = false
	return nil
}