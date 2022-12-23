package clock

import (
	"github.com/aalu1418/rpi-terminal-hub/services/base"
	"github.com/aalu1418/rpi-terminal-hub/types"
	log "github.com/sirupsen/logrus"
)

type service struct {
	types.Service
}

// provides a single handler for setting & incrementing metrics
func New(outgoingMsg chan<- types.Message) types.Service {
	var s service
	s.Service = base.New(outgoingMsg, types.CLOCK, types.CLOCK_FREQUENCY, s.onTick, s.processMsg)
	return &s
}

func (s *service) processMsg(m types.Message) {
	log.Warnf("[%s] received unexpected message: %s", s.Name(), m)
}

// poll & send message to metrics service
func (s *service) onTick() types.Message {
	msg := types.Message{
		To:   types.POSTOFFICE,
		Data: true,
	}

	return msg
}
