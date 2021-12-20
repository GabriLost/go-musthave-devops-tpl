package main

import (
	"flag"
	"github.com/GabriLost/go-musthave-devops-tpl/internal/agent"
	"github.com/GabriLost/go-musthave-devops-tpl/internal/types"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	pollInterval, reportIntervalFlag time.Duration
	addressFlag                      string
)

const (
	defaultAddress        = "localhost:8080"
	defaultPollInterval   = 2 * time.Second
	defaultReportInterval = 10 * time.Second
)

func getConfig() types.AgentConfig {
	var c types.AgentConfig
	err := env.Parse(&c)
	if err != nil {
		log.Println("Can't read env config")
		log.Println(err)
	}

	flag.StringVar(&addressFlag, "a", defaultAddress, "server address")
	flag.DurationVar(&reportIntervalFlag, "r", defaultReportInterval, "report interval")
	flag.DurationVar(&pollInterval, "p", defaultPollInterval, "poll interval")
	flag.Parse()

	//rewrite if ENV values is not empty
	_, isSet := os.LookupEnv("ADDRESS")
	if !isSet {
		c.Address = addressFlag
	}

	_, isSet = os.LookupEnv("REPORT_INTERVAL")
	if !isSet {
		c.ReportInterval = reportIntervalFlag
	}

	_, isSet = os.LookupEnv("POLL_INTERVAL")
	if !isSet {
		c.PollInterval = pollInterval
	}

	return c
}

func main() {
	//start processes
	types.SenderConfig = getConfig()
	types.SenderConfig.LogConfig()

	go agent.Schedule(agent.CollectRuntimeMetrics, types.SenderConfig.PollInterval)
	go agent.Schedule(agent.SendMetrics, types.SenderConfig.ReportInterval)

	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-cancelSignal

	os.Exit(1)
}
