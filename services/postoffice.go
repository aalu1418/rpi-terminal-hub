package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/aalu1418/rpi-terminal-hub/types"
	log "github.com/sirupsen/logrus"
)

// sorts the messages to the appropriate receiver like the post office
type postOffice struct {
	mailbag   <-chan types.Message
	mailboxes map[string]chan<- types.Message
	started   bool
	stop      chan struct{}
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
		m.mailboxes[address] = clients[i].ExtWrite()
	}

	return &m
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
			if msg.To == "" || msg.From == "" {
				log.Errorf("post office received invalid message: %+v", msg)
				continue
			}
			mailbox, exists := m.mailboxes[msg.To]
			if !exists {
				log.Errorf("no valid mailbox for addressee: %+v", msg)
				continue
			}
			mailbox <- msg
		}
	}
}

func (m *postOffice) Stop() error {
	if !m.started {
		return fmt.Errorf("post office has not been started")
	}
	close(m.stop)
	m.started = false
	return nil
}