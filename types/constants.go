package types

import "time"

const (
	// post office config
	POSTOFFICE  = "postoffice"
	MAX_WORKERS = 10 // max workers in post office processing messages

	// general service config
	MAX_QUEUE                     = 100       // max messages in queue
	INFINITE_TIME   time.Duration = 1<<63 - 1 // infinite time ~= 292.5 years
	DEFAULT_TIMEOUT               = 5 * time.Second

	// web server config
	WEBSERVER_ADDRESS = ":5000"
)
