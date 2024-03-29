//go:build !race

// skip race test - http handlers cannot be deregistered
package services

import (
	"context"
	"testing"

	"github.com/aalu1418/rpi-terminal-hub/services/alerts"
	"github.com/aalu1418/rpi-terminal-hub/services/base"
	"github.com/aalu1418/rpi-terminal-hub/services/connectivity"
	"github.com/aalu1418/rpi-terminal-hub/services/metrics"
	"github.com/aalu1418/rpi-terminal-hub/services/server"
	"github.com/aalu1418/rpi-terminal-hub/services/weather"
	"github.com/aalu1418/rpi-terminal-hub/types"
	"github.com/stretchr/testify/require"
)

// tests that all services must adhere to
func TestServices(t *testing.T) {
	deadline, _ := t.Deadline()
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	output := types.NewQueue()
	defer close(output)

	services := []types.Service{
		base.XXXNewBaseImplementation(t, output, "base", func(types.Message) {}),
		server.New(output),
		metrics.New(output),
		connectivity.New(output),
		weather.New(output, ""),
		alerts.NewNWS(output),
	}

	for _, v := range services {
		t.Run(v.Name(), func(t *testing.T) {
			t.Run("stopWithoutStart", func(t *testing.T) {
				require.Error(t, v.Stop(context.Background()))
			})
			t.Run("doubleStart", func(t *testing.T) {
				require.NoError(t, v.Start(ctx))
				require.Error(t, v.Start(ctx))
			})
			t.Run("channelClosure", func(t *testing.T) {
				write := v.ExtWrite()
				require.NoError(t, v.Stop(context.Background()))
				require.PanicsWithError(t, "send on closed channel", func() { write <- types.Message{} })
			})
		})
	}
}
