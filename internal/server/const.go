package server

const (
	MetricTypeGauge   = "gauge"
	MetricTypeCounter = "counter"

	contentType = "application/json"
)

var (
	metricCounters = make(map[string]int64)
	metricGauges   = make(map[string]float64)
)
