package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNextDuration(t *testing.T) {
	scheduleRaw := `{
		"Monday": "6:15PM",
		"Tuesday": "6:15PM",
		"Thursday": "6:15PM",
		"Friday": "6:15PM",
		"Saturday": "6:15PM"
	}`

	schedule := WeeklySchedule{}
	require.NoError(t, json.Unmarshal([]byte(scheduleRaw), &schedule))

	d := schedule.Next(time.Unix(0, 0))
	assert.Equal(t, 25*time.Hour+15*time.Minute, d)
}