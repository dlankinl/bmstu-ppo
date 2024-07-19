package web

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
	"time"
)

var requestMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "entrepreneurs",
	Subsystem:  "http",
	Name:       "request",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"handlerName", "status", "method"})

func observeRequest(d time.Duration, status int, method, handlerName string) {
	requestMetrics.WithLabelValues(handlerName, strconv.Itoa(status), method).Observe(d.Seconds())
}
