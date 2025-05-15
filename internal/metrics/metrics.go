package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	GRPCRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_server_requests_total",
			Help: "Total gRPC requests",
		},
		[]string{"method", "status"},
	)

	GRPCDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_server_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: []float64{0.1, 0.3, 1, 3, 5},
		},
		[]string{"method"},
	)
)

func Init() {
	prometheus.MustRegister(GRPCRequestsTotal, GRPCDuration)
}
