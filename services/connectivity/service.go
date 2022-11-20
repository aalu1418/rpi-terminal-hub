package connectivity

import (
	"net/http"

	"github.com/aalu1418/rpi-terminal-hub/services/base"
	"github.com/aalu1418/rpi-terminal-hub/types"
	log "github.com/sirupsen/logrus"
)

type service struct {
	types.Service
	client *http.Client
}

// provides a single handler for setting & incrementing metrics
func New(outgoingMsg chan<- types.Message) types.Service {
	var s service
	s.Service = base.New(outgoingMsg, types.CONNECTIVITY, types.CONN_FREQUENCY, s.onTick, s.processMsg)
	s.client = &http.Client{
		Timeout: types.CONN_TIMEOUT,
	}
	return &s
}

func (s *service) processMsg(m types.Message) {
	log.Warnf("[%s] received unexpected message: %s", s.Name(), m)
}

// poll & send message to metrics service
func (s *service) onTick() types.Message {
	msg := types.Message{
		To:   types.METRICS,
		Data: true,
	}

	if res, err := s.client.Get(types.CONN_URL); err != nil || res == nil || res.StatusCode != 200 {
		msg.Data = false
		log.Errorf("internet connectivity failed: %+v + error: %s", res, err)
	}

	return msg
}
