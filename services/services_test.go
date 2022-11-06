package services

import (
	"context"
	"testing"

	"github.com/aalu1418/rpi-terminal-hub/services/base"
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
	}

	for _, v := range services {
		t.Run(v.Name(), func(t *testing.T) {
			t.Run("stopWithoutStart", func(t *testing.T) {
				require.Error(t, v.Stop())
			})
			t.Run("doubleStart", func(t *testing.T) {
				require.NoError(t, v.Start(ctx))
				require.Error(t, v.Start(ctx))
			})
			t.Run("channelClosure", func(t *testing.T) {
				write := v.ExtWrite()
				require.NoError(t, v.Stop())
				require.PanicsWithError(t, "send on closed channel", func() { write <- types.Message{} })
			})
		})
	}
}
