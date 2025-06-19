package metrics

import (
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    RequestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "loadtester_requests_total",
        Help: "Total number of HTTP requests",
    })

    RequestsSuccess = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "loadtester_requests_successful",
        Help: "Total number of successful HTTP requests",
    })

    LatencyHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
        Name:    "loadtester_latency_seconds",
        Help:    "Histogram of response latencies",
        Buckets: prometheus.DefBuckets,
    })
)

func InitMetrics() {
    prometheus.MustRegister(RequestsTotal)
    prometheus.MustRegister(RequestsSuccess)
    prometheus.MustRegister(LatencyHistogram)
}

func StartPrometheusServer() {
    http.Handle("/metrics", promhttp.Handler())
    go http.ListenAndServe(":2112", nil)
}
