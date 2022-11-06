package services

import (
	"context"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aalu1418/rpi-terminal-hub/types"
	"github.com/aalu1418/rpi-terminal-hub/types/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMockService(t *testing.T, incomingChan chan<- types.Message) types.Service {
	s := mocks.NewService(t)
	s.On("Name").Return(fmt.Sprintf("mock-%d", rand.Int()))
	s.On("ExtWrite").Return(incomingChan)
	return s
}

func TestPostOffice(t *testing.T) {
	messages := types.NewQueue()
	defer close(messages)

	s0_msg := types.NewQueue()
	defer close(s0_msg)
	s0 := newMockService(t, s0_msg)
	s1_msg := types.NewQueue()
	defer close(s1_msg)
	s1 := newMockService(t, s1_msg)

	po := NewPostOffice(messages, []types.Service{s0, s1})

	assertEmptyQueues := func(t *testing.T) {
		for i := 0; i < 100; i++ {
			if len(messages) == 0 {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}

		assert.Equal(t, 0, len(messages), "main queue is not empty")
		assert.Equal(t, 0, len(s0_msg), "service 0 queue is not empty")
		assert.Equal(t, 0, len(s1_msg), "service 1 queue is not empty")
	}

	t.Run("start", func(t *testing.T) {
		require.Equal(t, types.POSTOFFICE, po.Name())
		require.NoError(t, po.Start(context.TODO()))
	})

	t.Run("happyPath", func(t *testing.T) {
		msg := types.Message{
			From: s0.Name(),
			To:   s1.Name(),
		}
		messages <- msg

		received := <-s1_msg
		assert.Equal(t, msg, received)
		assertEmptyQueues(t)
	})

	t.Run("missingTo", func(t *testing.T) {
		messages <- types.Message{
			From: s0.Name(),
		}
		assertEmptyQueues(t)
	})

	t.Run("missingFrom", func(t *testing.T) {
		messages <- types.Message{
			To: s1.Name(),
		}
		assertEmptyQueues(t)
	})

	t.Run("invalidTo", func(t *testing.T) {
		messages <- types.Message{
			From: s0.Name(),
			To:   "invalid",
		}
		assertEmptyQueues(t)
	})

	// multiple workers and dropping messages to handle too many messages
	t.Run("maxQueueMaxWorkers", func(t *testing.T) {
		totalMsgs := types.MAX_QUEUE * 2

		for i := 0; i < totalMsgs; i++ {
			messages <- types.Message{
				From: s0.Name(),
				To:   s1.Name(),
			}
		}

		for len(s1_msg) != types.MAX_QUEUE && len(messages) == 0 {
			// wait for queue to fill up, and incoming to finish processing
		}

		// cleanup messages
		var i atomic.Int32
	loop:
		for {
			select {
			case <-s1_msg:
				i.Store(0) // reset count if value foudn
			default:
				// break if 5 consecutive seconds of empty
				if i.Load() > 5 {
					break loop
				}
				time.Sleep(time.Second) // wait for 1 second to see if more txs show up
				i.Add(1)
			}
		}

		assertEmptyQueues(t)
	})

	t.Run("stop", func(t *testing.T) {
		require.NoError(t, po.Stop())

		// sent messages don't reach services
		messages <- types.Message{}

		assert.Equal(t, 1, len(messages))
		assert.Equal(t, 0, len(s0_msg))
		assert.Equal(t, 0, len(s1_msg))
	})
}