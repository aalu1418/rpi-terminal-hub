package infrared

import (
	"time"

	rpio "github.com/stianeikeland/go-rpio/v4"
)

func NewEmitter(pin rpio.Pin, pattern []int64) error {
	if err := rpio.Open(); err != nil {
		return err
	}
	defer rpio.Close()

	// convert to duration to save time during emit loop
	parsed := []time.Duration{}
	for i := range pattern {
		parsed = append(parsed, time.Duration(pattern[i]))
	}

	pin.Pwm()
	pin.DutyCycle(0, 2)     // set signal off
	pin.Freq(FREQUENCY * 2) // set clock frequency

	high := false
	for i := range parsed {
		s := time.Now()
		if !high {
			high = true
			pin.DutyCycle(1, 2) // pulse
		} else {
			high = false
			pin.DutyCycle(0, 2) // off
		}

		for time.Since(s) < parsed[i] {
			// pause
			// note: time.Sleep does not guarantee precision
		}
	}

	pin.DutyCycle(0, 2) // set back to low at the end
	return nil
}