package pricer

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	bidPrices *prometheus.GaugeVec
}

func newMetrics() *metrics {
	return &metrics{
		bidPrices: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "fastrpc",
			Subsystem: "pricer",
			Name:      "bid_price_gwei",
			Help:      "Bid price in gwei for different priority levels.",
		}, []string{"priority_level"}),
	}
}
