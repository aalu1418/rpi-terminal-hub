package signals

import (
	"fmt"
	"testing"
	"time"
)

func TestSignal(t *testing.T) {
	s := Signal{
		Data: []time.Duration{
			100 * time.Nanosecond,
			100 * time.Nanosecond,
			100 * time.Nanosecond,
			100 * time.Nanosecond,
			100 * time.Nanosecond,
			100 * time.Nanosecond,
		},
	}

	tick := time.Now()
	logTime := func() {
		fmt.Println(time.Since(tick))
		tick = time.Now()
	}

	s.Run(logTime, logTime)
	time.Sleep(time.Second)
}