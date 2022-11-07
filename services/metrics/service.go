package metrics

import (
	"fmt"

	"github.com/aalu1418/rpi-terminal-hub/services/base"
	"github.com/aalu1418/rpi-terminal-hub/types"
	log "github.com/sirupsen/logrus"
)

const (
	NAME = "prometheus-metrics"
)

type service struct {
	types.Service
}

// provides a single handler for setting & incrementing metrics
func New(outgoingMsg chan<- types.Message) types.Service {
	var s service
	s.Service = base.New(outgoingMsg, NAME, types.INFINITE_TIME, s.onTick, s.processMsg)
	return &s
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
