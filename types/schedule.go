package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func StringToDay(s string) (time.Weekday, error) {
	for i := 0; i < 7; i++ {
		if strings.ToLower(time.Weekday(i).String()) == strings.ToLower(s) {
			return time.Weekday(i), nil
		}
	}
	return 0, fmt.Errorf("invalid day string: %s", s)
}

type WeeklySchedule map[time.Weekday]time.Time

func (w WeeklySchedule) UnmarshalJSON(b []byte) error {
	raw := map[string]string{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	for k, v := range raw {
		parsedTime, err := time.Parse("3:04PM", v)
		if err != nil {
			return err
		}
		day, err := StringToDay(k)
		if err != nil {
			return err
		}

		w[day] = parsedTime
	}
	return nil
}

func (w WeeklySchedule) Next(current time.Time) time.Duration {
	day := current.Weekday()

	var next time.Time
	days := 0
	for {
		t, exists := w[day]
		if exists {
			// build next time relative to current time
			y, m, d := current.AddDate(0, 0, days).Date()                             // get date with added days
			next = time.Date(y, m, d, t.Hour(), t.Minute(), 0, 0, current.Location()) // build time with info

			if next.After(current) {
				break
			}
		}
		days++
		day = (day + 1) % 7
	}

	return next.Sub(current) // calculate duration until next
}