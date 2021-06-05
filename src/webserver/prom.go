// promethues

package webserver

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// prometheus namespace
	promNamespace = "webserver"
	// gin prometheus labels
	promLabels = []string{
		"status_code",
		"path",
		"method",
	}
	promUptime = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: promNamespace,
			Name:      "server_uptime",
			Help:      "gin server uptime in seconds",
		}, nil,
	)
	promReqCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: promNamespace,
			Name:      "req_count",
			Help:      "gin server request count",
		}, promLabels,
	)
	promReqLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: promNamespace,
			Name:      "req_latency",
			Help:      "gin server request latency in seconds",
		}, promLabels,
	)
)

// PromExporterHandler return a handler as the prometheus metrics exporter
func PromExporterHandler(collectors ...prometheus.Collector) gin.HandlerFunc {
	// uptime
	go func() {
		for range time.Tick(time.Second) {
			promUptime.WithLabelValues().Inc()
		}
	}()
	return gin.WrapH(promhttp.Handler())
}
