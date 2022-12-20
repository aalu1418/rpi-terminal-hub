package types

import (
	"strings"
	"time"
)

// only defines a subset of interested parameters
type NWSAlerts struct {
	Title    string     `json:"title"`
	Updated  time.Time  `json:"updated"`
	Features []NWSAlert `json:"features"`
}

type NWSAlert struct {
	Properties NWSAlertProps
}

type NWSAlertProps struct {
	Event       string    `json:"event"`
	Headline    string    `json:"headline"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	Effective   time.Time `json:"effective"`
	Expires     time.Time `json:"expires"`
	Ends        time.Time `json:"ends"`
}

const (
	ADVISORY = AlertLevel(iota)
	WATCH
	WARNING
	INVALID
)

var (
	alertLevelStr = []string{"Advisory", "Watch", "Warning"}
)

type AlertLevel uint8

func (al AlertLevel) String() string {

	if int(al) >= len(alertLevelStr) {
		return "NA"
	}

	return alertLevelStr[al]
}

func ParseAlertLevel(s string) AlertLevel {
	for i, v := range alertLevelStr {
		if strings.EqualFold(v, s) {
			return AlertLevel(i)
		}
	}
	return INVALID
}

type ParsedAlert struct {
	Level AlertLevel
	Start time.Time
	End   time.Time
}