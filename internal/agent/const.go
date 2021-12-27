package agent

import "time"

const (
	DefaultProtocol = "http://"
	TCP             = "tcp"
	DefaultTimeout  = 2 * time.Second
	DefaultAddress  = "localhost:8080"
)

var Metrics []Gauge
var PollCount int64
