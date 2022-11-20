package vacuum

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aalu1418/rpi-terminal-hub/rpio/infrared"
	"github.com/aalu1418/rpi-terminal-hub/services/base"
	"github.com/aalu1418/rpi-terminal-hub/types"
	log "github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
)

var Calls = []string{"start", "stop", "30min", "home"}

type service struct {
	types.Service
	once     bool
	schedule types.WeeklySchedule
	stop     chan struct{}
}

// provides a single handler for setting & incrementing metrics
func New(outgoingMsg chan<- types.Message) types.Service {
	var s service
	// 1s tick - continuosly run
	s.Service = base.New(outgoingMsg, types.VACUUM, time.Second, s.onTick, s.processMsg)
	s.schedule = types.WeeklySchedule{}
	s.stop = make(chan struct{})
	return &s
}

func (s *service) Start(ctx context.Context) error {
	if err := json.Unmarshal([]byte(types.VACUUM_SCHEDULE), &s.schedule); err != nil {
		return err
	}

	return s.Service.Start(ctx)
}

func (s *service) Stop(ctx context.Context) error {
	close(s.stop)
	return s.Service.Stop(ctx)
}

func (s *service) processMsg(m types.Message) {
	if m.To != types.VACUUM || m.From != types.WEBSERVER {
		log.Errorf("[VACUUM] received incorrect message: %+v", m)
		return
	}

	command := strings.ToLower(m.Data.(string))
	switch command {
	case "start":
		fallthrough
	case "stop":
		if err := infrared.NewEmitter(rpio.Pin(types.IR_EMITTER), COMMAND_STARTSTOP); err != nil {
			log.Errorf("[VACUUM] failed to emit START/STOP command: %s", err)
		}
	case "30min":
		if err := infrared.NewEmitter(rpio.Pin(types.IR_EMITTER), COMMAND_30MIN); err != nil {
			log.Errorf("[VACUUM] failed to emit 30MIN command: %s", err)
		}
	case "home":
		if err := infrared.NewEmitter(rpio.Pin(types.IR_EMITTER), COMMAND_HOME); err != nil {
			log.Errorf("[VACUUM] failed to emit HOME command: %s", err)
		}
	default:
		log.Errorf("[VACUUM] unknown command: %s", command)
	}
}

// called once at the very beginning (and after INFINITE_TIME)
func (s *service) onTick() types.Message {
	if !s.once {
		s.once = true
		return types.Message{
			To:   types.POSTOFFICE,
			Data: "[ALIVE] " + s.Name(),
		}
	}

	// calculate next trigger time
	d := s.schedule.Next(time.Now())
	tick := time.After(d)
	log.Infof("[VACUUM] scheduled in %s", d)
	select {
	case <-s.stop:
		fmt.Println("interrupted")
		return types.Message{
			To:   types.POSTOFFICE,
			Data: "[VACUUM] shutting down trigger",
		}
	case <-tick:
		// continue
	}

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
