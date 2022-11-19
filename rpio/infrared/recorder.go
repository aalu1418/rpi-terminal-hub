package infrared

import (
	"fmt"
	"log"
	"strings"
	"time"

	rpio "github.com/stianeikeland/go-rpio/v4"
)

func NewRecorder(pin rpio.Pin) error {
	if err := rpio.Open(); err != nil {
		return err
	}
	defer rpio.Close()

	pin.Input()

	log.Print("[IR] recorder ready")
	var start = time.Now()
	var detected bool
	var timings []time.Duration
	state := pin.Read()
	for {
		since := time.Since(start)
		newState := pin.Read()
		if state != newState { // if edge changes
			detected = true
			timings = append(timings, since) // log how long edge was in current state
			start = time.Now()               // save change timestamp
			state = newState
			continue // continue loop with evaluating exit condition
		}

		// exit if signal detected and there is a 1s gap (likely signal is ended)
		// exit if no signal detected and there is 10s waiting
		if (detected && since > time.Second) || (!detected && since > 10*time.Second) {
			break
		}
	}
	if !detected {
		log.Print("[IR] no signal")
		return nil
	}

	output := []string{}
	for _, v := range timings {
		output = append(output, fmt.Sprintf("%d", int64(v)))
	}
	output = output[1:]

	log.Printf("[]int64{" + strings.Join(output, ",") + "}")

	return nil
}
