package agent

import (
	"sync"
	"time"
)

const (
	DefaultProtocol = "http://"
	TCP             = "tcp"
	DefaultTimeout  = 2 * time.Second
	DefaultAddress  = "localhost:8080"
)

var Metrics []Gauge
var PollCount int64

type utilizationData struct {
	mu              sync.Mutex
	TotalMemory     Gauge
	FreeMemory      Gauge
	CPUutilizations []Gauge
	CPUtime         []float64
	CPUutilLastTime time.Time
}

var UtilizationData utilizationData
