package signals

import (
	"time"
)

type Signal []time.Duration

func (s Signal) Run(high, low func()) {
	for i, v := range s {
		start := time.Now()
		if i%2 == 0 { // every other starting with setting high
			high()
		} else {
			low()
		}

		// tests indicate mostly ok for 0.5 microsecond precision +/- 10%
		for time.Since(start) < v {
			// while loop
		}
	}

	// ensure low after finish run
	low()
}

// for creating a digital PWM signal
func (s Signal) Freq(hz int) (new Signal) {
	period := time.Second / time.Duration(hz)

	for i, v := range s {
		if i%2 == 0 { // high
			for j := 0; j < int(v/period); j++ {
				new = append(new, period/2)
				new = append(new, period/2)
			}

		} else { // low
			new = append(new, v)
		}
	}

	return new
}
