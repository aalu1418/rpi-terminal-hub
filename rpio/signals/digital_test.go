package signals

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignal(t *testing.T) {
	d := []time.Duration{}
	N := 10000
	inc := 500 * time.Nanosecond

	s := Signal{}
	for i := 0; i < N; i++ {
		s = append(s, inc)
	}

	tick := time.Now()
	logTime := func() {
		d = append(d, time.Since(tick))
		tick = time.Now()
	}

	s.Run(logTime, logTime)
	time.Sleep(time.Duration(N*150/100) * inc) // wait for completion

	// skip first count (time for set up)
	var sum time.Duration
	for _, v := range d[1:] {
		sum += v
	}

	// calc average
	avg := sum / time.Duration(len(d)-1)

	// calc percent difference
	diff := (avg - inc).Abs()

	// assert 10% > difference
	assert.GreaterOrEqual(t, inc*1/10, diff)
	t.Logf("Avg %s", avg)
}

func TestSignal_Freq(t *testing.T) {
	freq := 10 // 10 hz
	halfPeriod := time.Second / time.Duration(freq*2)

	signal := Signal{time.Second, time.Second, time.Second}
	pwm := signal.Freq(freq)

	var expected Signal
	// first PWM burst
	for i := 0; i < int(time.Second/halfPeriod); i++ {
		expected = append(expected, halfPeriod)
	}
	// signal low
	expected = append(expected, time.Second)
	// second PWM burst
	for i := 0; i < int(time.Second/halfPeriod); i++ {
		expected = append(expected, halfPeriod)
	}

	assert.Equal(t, expected, pwm)
}