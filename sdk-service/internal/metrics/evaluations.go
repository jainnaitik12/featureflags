package metrics

import "github.com/prometheus/client_golang/prometheus"

// FlagEvaluations counts flag reads (per flag_name label).
var FlagEvaluations = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "flag_evaluations_total",
		Help: "Total number of flag evaluations",
	},
	[]string{"flag_name"},
)

// Register adds Prometheus collectors for this service.
func Register() {
	prometheus.MustRegister(FlagEvaluations)
}
