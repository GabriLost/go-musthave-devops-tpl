package main

import (
	"github.com/GabriLost/go-musthave-devops-tpl/internal/agent"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	pollRuntimeMetricsInterval = 2 * time.Second
	reportInterval             = 10 * time.Second
)

func main() {
	go agent.Schedule(agent.StoreRuntimeMetrics, pollRuntimeMetricsInterval)
	go agent.Schedule(agent.SendMetrics, reportInterval)

	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-cancelSignal

	os.Exit(1)
}
