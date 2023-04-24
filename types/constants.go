package types

import "time"

// non-configurable parameters
const (
	// post office config
	POSTOFFICE = "postoffice"

	// general service config
	MAX_QUEUE                   = 100       // max messages in queue
	INFINITE_TIME time.Duration = 1<<63 - 1 // infinite time ~= 292.5 years

	// web server config
	WEBSERVER = "web-server"

	// metrics config
	METRICS = "metrics-handler"

	// connectivity config
	CONNECTIVITY = "internet-connectivity"

	// weather config
	WEATHER    = "weather"
	OWM_ENVVAR = "OWM"

	// NWS alert config
	NWS = "alert-nws"

	// vacuum config
	VACUUM = "vacuum"
)

// configurable parameters via -ldflags
// go run -ldflags "-X github.com/aalu1418/rpi-terminal-hub/types.WEATHER_LOCATION=hello" main.go
var (
	// post office config
	MAX_WORKERS int32 = 10 // max workers in post office processing messages

	// general service config
	DEFAULT_TIMEOUT = 5 * time.Second

	// web server config
	WEBSERVER_ADDRESS = "0.0.0.0:5000"

	// connectivity config
	CONN_FREQUENCY = 5 * time.Second  // frequency of polling
	CONN_TIMEOUT   = 1 * time.Second  // timeout for request
	CONN_URL       = "http://1.1.1.1" // endpoint to ping

	// weather config
	WEATHER_FREQUENCY = 15 * time.Minute
	WEATHER_LOCATION  = "Westminster, CO"
	WEATHER_LAT       = "39.86506214649686"
	WEATHER_LON       = "-105.04846274505923"

	// NWS alert config
	NWS_FREQUENCY = 5 * time.Minute

	// IR config
	IR_RECEIVER = 22
	IR_EMITTER  = 18

	// vacuum config
	VACUUM_SCHEDULE = `{
	  "Sunday": "6:15PM",
	  "Tuesday": "6:15PM",
	  "Wednesday": "6:15PM",
	  "Thursday": "6:15PM",
	  "Saturday": "6:15PM"
	}`
)
