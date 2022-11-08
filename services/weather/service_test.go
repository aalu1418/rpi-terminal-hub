package weather

import (
	"context"
	"testing"

	"github.com/aalu1418/rpi-terminal-hub/types"
	"github.com/stretchr/testify/require"
)

func TestWeather(t *testing.T) {
	ctx := context.Background()

	msgs := types.NewQueue()
	defer close(msgs)

	s := New(msgs)

	require.NoError(t, s.Start(ctx))

	<-msgs

	require.NoError(t, s.Stop(ctx))
}