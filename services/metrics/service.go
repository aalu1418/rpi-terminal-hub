package metrics

import (
	"fmt"

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
	s.Service = base.New(outgoingMsg, types.METRICS, types.INFINITE_TIME, s.onTick, s.processMsg)
	return &s
}

func (s *service) processMsg(m types.Message) {
	var ok bool
	switch m.From {
	case types.CONNECTIVITY:
		var success bool
		if success, ok = m.Data.(bool); !ok {
			break
		}
		if success {
			internetAlive.Inc()
		}
	default:
		log.Warnf("[%s] unexpected sender: %+v", s.Name(), m)
		return
	}

	if !ok {
		log.Errorf("[%s] invalid data type: %+v", s.Name(), m)
	}
}

// called once at the very beginning (and after INFINITE_TIME)
func (s *service) onTick() types.Message {
	return types.Message{
		To:   types.POSTOFFICE,
		Data: fmt.Sprintf("[ALIVE] %s", s.Name()),
	}
}
