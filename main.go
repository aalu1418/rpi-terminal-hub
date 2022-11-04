package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	system := make(chan os.Signal, 1)
	signal.Notify(system, os.Interrupt)

	// services := []types.Service{}

	// central message routing loop
	for {
		select {
		case <-system:
			fmt.Println("closing!")
			return
		default:
			fmt.Println("looping")
		}
		time.Sleep(1 * time.Second)
	}
}
