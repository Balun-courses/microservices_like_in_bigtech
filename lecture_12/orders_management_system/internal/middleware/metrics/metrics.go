package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	appName   = "orders_management_system"
	namespace = appName
)

var ms struct {
	responseTimeHistogram *prometheus.HistogramVec
}

func init() {
	ms.responseTimeHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "grpc",
			Name:      appName + "_histogram_response_time_seconds",
			Help:      "Время ответа от сервера",
			Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
		},
		[]string{"method", "is_error"},
	)
}

func responseTimeHistogramObserve(method string, err error, d time.Duration) {
	isError := strconv.FormatBool(err != nil)
	ms.responseTimeHistogram.WithLabelValues(method, isError).Observe(d.Seconds())
}
