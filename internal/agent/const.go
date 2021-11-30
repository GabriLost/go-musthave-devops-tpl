package agent

import "time"

const (
	DefaultProtocol = "http://"
	DefaultServer   = "localhost"
	DefaultPort     = "8080"
	ContentType     = "text/plain"
	TCP             = "tcp"
	DefaultTimeout  = 2 * time.Second
)

var Metrics []Gauge
var PollCount int64
