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
	pollInterval, reportInterval int
	address                      string
)

func main() {
	var cfg types.AgentConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Println("Can't read env config")
		log.Println(err)
	}

	//rewrite if flags is not empty
	flag.StringVar(&address, "a", "", "server address")
	flag.IntVar(&reportInterval, "r", 0, "report interval")
	flag.IntVar(&pollInterval, "p", 0, "poll interval")
	flag.Parse()
	if address != "" {
		types.SenderConfig.Address = address
	}
	if pollInterval != 0 {
		types.SenderConfig.PollInterval = time.Second * time.Duration(pollInterval)
	}
	if reportInterval != 0 {
		types.SenderConfig.ReportInterval = time.Second * time.Duration(reportInterval)
	}

	//start processes
	go agent.Schedule(agent.CollectRuntimeMetrics, types.SenderConfig.PollInterval)
	go agent.Schedule(agent.SendMetrics, types.SenderConfig.ReportInterval)

	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-cancelSignal

	os.Exit(1)
}
