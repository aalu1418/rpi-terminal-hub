package vacuum

import (
	"fmt"
	"time"

	"github.com/aalu1418/rpi-terminal-hub/rpio/infrared"
	"github.com/aalu1418/rpi-terminal-hub/services/base"
	"github.com/aalu1418/rpi-terminal-hub/types"
	"github.com/stianeikeland/go-rpio/v4"
)

type service struct {
	types.Service
	notFirst bool
}

// provides a single handler for setting & incrementing metrics
func New(outgoingMsg chan<- types.Message) types.Service {
	var s service
	// 1s tick - continuosly run
	s.Service = base.New(outgoingMsg, types.VACUUM, time.Second, s.onTick, s.processMsg)
	return &s
}

func (s *service) processMsg(m types.Message) {
}

// called once at the very beginning (and after INFINITE_TIME)
func (s *service) onTick() types.Message {
	if s.notFirst {
		s.notFirst = true
		return types.Message{
			To:   types.POSTOFFICE,
			Data: "[ALIVE] " + s.Name(),
		}
	}

	// TODO:calculate next trigger time
	tick := time.After(0)

	<-tick
	if err := infrared.NewEmitter(rpio.Pin(types.IR_EMITTER), COMMAND_30MIN); err != nil {
		return types.Message{
			To:   types.POSTOFFICE,
			Data: fmt.Sprintf("[ERROR] vacuum failed to emit signal: %s", err),
		}
	}

	return types.Message{
		To:   types.POSTOFFICE,
		Data: "[VACUUM] IR signal emitted",
	}
}
