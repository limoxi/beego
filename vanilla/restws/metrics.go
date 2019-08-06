package restws

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
)

var restwsGauge = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "restws_connection",
	Help: "Number of restws connection is active",
})

func GetRestwsGauge() prometheus.Gauge {
	return restwsGauge
}

func init() {
}
