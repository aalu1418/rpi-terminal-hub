package ws2812b

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"github.com/aalu1418/rpi-terminal-hub/rpio/signals"
	rpio "github.com/stianeikeland/go-rpio/v4"
)

// https://www.arrow.com/en/research-and-events/articles/protocol-for-the-ws2812b-programmable-led

var (
	ONE_HIGH  = 650 * time.Nanosecond
	ONE_LOW   = 300 * time.Nanosecond
	ZERO_HIGH = 250 * time.Nanosecond
	ZERO_LOW  = 700 * time.Nanosecond
	RESET     = 300 * time.Microsecond
)

type LEDs []color.Color

func New(pin rpio.Pin, leds LEDs) error {
	if err := rpio.Open(); err != nil {
		return err
	}
	defer rpio.Close()

	s := leds.Build()
	fmt.Println(s)

	pin.Output()

	s.Run(
		func() { pin.High() },
		func() { pin.Low() },
	)

	return nil
}

func (l LEDs) Build() (s signals.Signal) {
	for _, v := range l {
		r, g, b, _ := v.RGBA()

		// GRB encoding
		s = append(s, encode(g)...)
		s = append(s, encode(r)...)
		s = append(s, encode(b)...)
	}

	s[len(s)-1] += RESET // add reset to last LOW signal
	return s
}

func encode(c uint32) (s signals.Signal) {
	str := fmt.Sprintf("%08b", uint8(c))

	for _, v := range strings.Split(str, "") {
		if v == "0" {
			s = append(s, ZERO_HIGH)
			s = append(s, ZERO_LOW)
		} else {
			s = append(s, ONE_HIGH)
			s = append(s, ONE_LOW)
		}
	}

	return s
}