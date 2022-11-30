package signals

import (
	"time"
)

type Signal struct {
	Data []time.Duration
}

func (s Signal) Run(high, low func()) {
	for i, v := range s.Data {
		start := time.Now()
		if i%2 == 0 { // every other starting with setting high
			high()
		} else {
			low()
		}

		for time.Since(start) < v {
			// while loop
		}
	}
}
