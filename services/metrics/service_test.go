package metrics

import (
	"testing"

	"github.com/aalu1418/rpi-terminal-hub/types"
	"github.com/aalu1418/rpi-terminal-hub/types/mocks"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetrics_ProcessMsg(t *testing.T) {
	var s service
	base := mocks.NewService(t)
	base.On("Name").Return(t.Name())
	s.Service = base

	t.Run("unexpectedFrom", func(t *testing.T) {
		s.processMsg(types.Message{
			From: "unexpected",
		})
		require.Equal(t, float64(0), testutil.ToFloat64(internetAlive))
	})

	t.Run(types.CONNECTIVITY, func(t *testing.T) {
		init := testutil.ToFloat64(internetAlive)
		baseMsg := types.Message{
			From: types.CONNECTIVITY,
		}

		data := []struct {
			name   string
			data   interface{}
			change bool
		}{
			{"invalid", "invalid", false},
			{"false", false, false},
			{"true", true, true},
		}

		for _, d := range data {
			t.Run(d.name, func(t *testing.T) {
				msg := baseMsg
				msg.Data = d.data
				s.processMsg(msg)

				if d.change {
					assert.Equal(t, init+1, testutil.ToFloat64(internetAlive))
				} else {
					assert.Equal(t, init, testutil.ToFloat64(internetAlive))
				}
			})
		}
	})
}