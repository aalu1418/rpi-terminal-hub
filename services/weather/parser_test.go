package weather

import (
	"testing"

	"github.com/aalu1418/rpi-terminal-hub/types"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func TestParser(t *testing.T) {
	N := 4
	// test handling if hourly length < 12, wrap around to not panic
	input := types.OneCall{}
	for i := 0; i < 4; i++ {
		input.Hourly = append(input.Hourly, types.Weather{
			Temp: float64(i), // use temp as increment tracking
		})
	}

	output := Parser(input)
	for i, h := range output.Hourly {
		assert.Equal(t, (3*i+1)%N, h.Temperature)
	}
}

func TestWeatherParser(t *testing.T) {
	assert.NotPanics(t, func() {
		weatherParser(types.Weather{})
	})
}

func TestIconParser(t *testing.T) {
	inputs := []string{}
	keys := maps.Keys(types.IconMap)
	for i := range keys {
		if len(keys[i]) != 3 {
			inputs = append(inputs, keys[i]+"d")
			inputs = append(inputs, keys[i]+"n")
		} else {
			inputs = append(inputs, keys[i])
		}
	}

	// test expected inputs
	for _, v := range inputs {
		t.Run(v, func(t *testing.T) {
			assert.NotEqual(t, types.IconNA, iconParser(v))
		})
	}

	invalidInputs := []string{
		"", "01", "04", "0", "blah",
	}
	for _, v := range invalidInputs {
		t.Run(v, func(t *testing.T) {
			assert.Equal(t, types.IconNA, iconParser(v))
		})
	}

}