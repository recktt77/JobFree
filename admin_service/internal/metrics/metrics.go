package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var BanUserCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "admin_ban_user_total",
		Help: "Total number of banned users",
	},
)

func Register() {
	prometheus.MustRegister(BanUserCounter)
}

func Handler() http.Handler {
	return promhttp.Handler()
}
