package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
)

var endpointCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "endpoint_call_total",
		Help: "total counts for panic",
	},
	[]string{"endpoint", "method"},
)

/*
var endpointSummary = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "endpoint_durations_seconds",
		Help:       "endpoint latency distributions.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"endpoint", "method"},
)*/

//var normDomain = 0.0002
//var normMean = 0.00001
var endpointHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "endpoint_durations_histogram_seconds",
		Help:    "endpoint latency distributions.",
	//Buckets: prometheus.LinearBuckets(normMean-5*normDomain, .5*normDomain, 20),
	},
	[]string{"endpoint", "method"},
)

var panicCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "panic_total",
	Help: "total counts for panic",
})

var resourceRetryCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "resource_retry_total",
	Help: "total counts for resource's retry",
})

var businessErrorCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "business_error_total",
	Help: "total counts for business error",
})

var restwsGauge = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "restws_connection",
	Help: "Number of rest proxy websocket connection is active",
})

func GetRestwsGauge() prometheus.Gauge {
	return restwsGauge
}

func GetEndpointCounter() *prometheus.CounterVec {
	return endpointCounter
}

func GetEndpointSummary() *prometheus.SummaryVec {
	return nil
}

func GetEndpointHistogram() *prometheus.HistogramVec {
	return endpointHistogram
}

func GetPanicCounter() prometheus.Counter {
	return panicCounter
}

func GetBusinessErrorCounter() prometheus.Counter {
	return businessErrorCounter
}

func GetResourceRetryCounter() prometheus.Counter {
	return resourceRetryCounter
}

func init() {
	fmt.Println("in metrics init")
}