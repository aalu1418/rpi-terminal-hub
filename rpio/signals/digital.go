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

	var next time.Duration
	for i, v := range s {
		if i%2 == 0 { // high
			count := int(v / period)
			for j := 0; j < count; j++ {
				new = append(new, period/2) // high

				// handle cases where last period
				if j == count-1 {
					// if last signal, append low
					if i == len(s)-1 {
						new = append(new, period/2)
					} else {
						// if last period of pulse, combine with low later
						next = period / 2
					}

				} else {
					new = append(new, period/2)
				}
			}

		} else { // low
			new = append(new, v+next)
			next = 0
		}
	}

	return new
}
