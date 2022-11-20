package types

import (
	"encoding/json"
	"strconv"
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
	tz := time.Now().Format("-07")
	shift, err := strconv.Atoi(tz)
	require.NoError(t, err)
	assert.Equal(t, 18*time.Hour+15*time.Minute-time.Duration(shift)*time.Hour, d)
}