package base

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aalu1418/rpi-terminal-hub/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBase(t *testing.T) {
	name := "test"
	msg := types.Message{
		From: "external",
		To:   "test",
		Data: time.Now(),
	}

	var wg sync.WaitGroup
	wg.Add(1) // should only be called once
	processMsg := func(m types.Message) {
		assert.Equal(t, msg, m)
		wg.Done()
	}

	s := XXXNewBaseImplementation(t, strings.ToUpper(name), processMsg)

	// start without issue
	assert.Equal(t, name, s.Name())
	require.NoError(t, s.Start(context.Background()))
	assert.True(t, s.Healthy())

	// check received message
	read := s.ExtRead()
	receivedMsg := <-read
	assert.Equal(t, name, receivedMsg.From)
	assert.Equal(t, "placeholder", receivedMsg.To)
	assert.Equal(t, "time.Time", fmt.Sprintf("%T", receivedMsg.Data))

	write := s.ExtWrite()
	write <- types.Message{To: name}                            // invalid
	write <- types.Message{From: "placeholder", To: "not_test"} // invalid
	write <- msg
	wg.Wait()

	require.NoError(t, s.Stop())
	assert.False(t, s.Healthy())
}