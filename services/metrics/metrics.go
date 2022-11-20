package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	internetAlive = promauto.NewCounter(prometheus.CounterOpts{
		Name: "internet_alive",
		Help: "The total number of successful internet pings",
	})
)