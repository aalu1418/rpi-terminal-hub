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
	WEBSERVER         = "web-server"
	WEBSERVER_ADDRESS = "0.0.0.0:5000"

	// metrics config
	METRICS = "metrics-handler"

	// connectivity config
	CONNECTIVITY   = "internet-connectivity"
	CONN_FREQUENCY = 5 * time.Second  // frequency of polling
	CONN_TIMEOUT   = 1 * time.Second  // timeout for request
	CONN_URL       = "http://1.1.1.1" // endpoint to ping

	// weather config
	WEATHER           = "weather"
	WEATHER_FREQUENCY = 15 * time.Minute
	WEATHER_LAT       = "39.86506214649686"
	WEATHER_LON       = "-105.04846274505923"
)
