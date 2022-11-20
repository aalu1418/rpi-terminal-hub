package types

import "time"

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
}