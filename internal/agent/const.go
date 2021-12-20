package agent

import "time"

const (
	DefaultProtocol = "http://"
	TCP             = "tcp"
	DefaultTimeout  = 2 * time.Second
)

var Metrics []Gauge
var PollCount int64
