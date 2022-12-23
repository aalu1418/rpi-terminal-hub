package clock

import (
	"image/color"
	"time"

	"github.com/aalu1418/rpi-terminal-hub/rpio/ws2812b"
	"github.com/stianeikeland/go-rpio/v4"
)

type Off struct{}

func (o Off) RGBA() (r, g, b, a uint32) {
	return 125, 125, 0, 0
}

func Demo(pin rpio.Pin) error {

	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			if err := ws2812b.New(pin, []color.Color{color.White}); err != nil {
				return err
			}

		} else {
			if err := ws2812b.New(pin, []color.Color{color.Black}); err != nil {
				return err
			}

		}

		time.Sleep(time.Second)
	}

	return ws2812b.New(pin, []color.Color{color.Black})
}