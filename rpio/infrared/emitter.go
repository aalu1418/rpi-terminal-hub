package infrared

import (
	"time"

	"github.com/aalu1418/rpi-terminal-hub/rpio/signals"
	rpio "github.com/stianeikeland/go-rpio/v4"
)

func NewEmitter(pin rpio.Pin, pattern []int64) error {
	if err := rpio.Open(); err != nil {
		return err
	}
	defer rpio.Close()

	// convert to duration to save time during emit loop
	parsed := signals.Signal{}
	for i := range pattern {
		parsed = append(parsed, time.Duration(pattern[i]))
	}

	pin.Pwm()
	pin.DutyCycle(0, 2)     // set signal off
	pin.Freq(FREQUENCY * 2) // set clock frequency

	parsed.Run(
		func() { pin.DutyCycle(1, 2) }, // high
		func() { pin.DutyCycle(0, 2) }, // low
	)
	return nil
}