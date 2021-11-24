package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

//метрика gauge
type gauge struct {
	name  string
	value float64
}

//метрика counter
type counter struct {
	name  string
	value int64
}

const (
	defaultServer  = "http://localhost"
	defaultPort    = "8080"
	contentType    = "text/plain"
	defaultTimeout = 2 * time.Second
)

const (
	pollRuntimeMetricsInterval = 2 * time.Second
	reportInterval             = 10 * time.Second
)

var Metrics []gauge
var PollCount int64

// SendGauge отправить данные на сервер
func (g gauge) SendGauge(client *http.Client) (bool, error) {
	url := fmt.Sprintf(defaultServer+":"+defaultPort+"/update/%s/%s/%d", "gauge", g.name, int(g.value))
	log.Println("SendGauge " + url)
	resp, err := client.Post(url, "text/plain", nil)
	if err != nil {
		log.Println(err)
		return false, err
	}
	err = resp.Body.Close()
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}

func (c counter) SendCounter(client *http.Client) (bool, error) {
	url := fmt.Sprintf(defaultServer+":"+defaultPort+"/update/%s/%s/%d", "counter", c.name, c.value)
	log.Println("SendCounter " + url)
	resp, err := client.Post(url, contentType, nil)
	if err != nil {
		log.Println(err)
		return false, err
	}
	err = resp.Body.Close()
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}

func SendDataAsync() {
	ticker := time.NewTicker(reportInterval)
	for range ticker.C {
		SendData(Metrics)
	}
}

func SendData(metrics []gauge) {
	log.Println("Total metrics is ", len(metrics))
	metrics = AddRandomMetrics(metrics)
	log.Println("Total metrics is ", len(metrics))
	client := http.Client{Timeout: defaultTimeout}
	for _, i := range metrics {
		_, err := i.SendGauge(&client)
		if err != nil {
			return
		}
	}
	metricCounter := counter{name: "PollCount", value: PollCount}
	log.Println("Reset poll counter to zero")
	PollCount = 0
	metricCounter.SendCounter(&client)
	client.CloseIdleConnections()
}

//AddRandomMetrics метрипки не относятся к рантайму, добавляем отдельно
func AddRandomMetrics(metrics []gauge) []gauge {
	randomMetric := gauge{name: "RandomValue", value: rand.Float64() * 100}
	metrics = append(metrics, randomMetric)
	return metrics
}

func GetRuntimeMetrics() {
	var rtm runtime.MemStats
	log.Println("ticker GetRuntimeMetrics")
	runtime.ReadMemStats(&rtm)
	PollCount += 1
	Metrics = []gauge{
		{name: "Alloc", value: float64(rtm.Alloc)},
		{name: "BuckHashSys", value: float64(rtm.BuckHashSys)},
		{name: "Frees", value: float64(rtm.Frees)},
		{name: "GCCPUFraction", value: rtm.GCCPUFraction},
		{name: "GCSys", value: float64(rtm.GCSys)},
		{name: "HeapAlloc", value: float64(rtm.HeapAlloc)},
		{name: "HeapIdle", value: float64(rtm.HeapIdle)},
		{name: "HeapInuse", value: float64(rtm.HeapInuse)},
		{name: "HeapObjects", value: float64(rtm.HeapObjects)},
		{name: "HeapReleased", value: float64(rtm.HeapReleased)},
		{name: "HeapSys", value: float64(rtm.HeapSys)},
		{name: "LastGC", value: float64(rtm.LastGC)},
		{name: "Lookups", value: float64(rtm.Lookups)},
		{name: "MCacheInuse", value: float64(rtm.MCacheInuse)},
		{name: "MCacheSys", value: float64(rtm.MCacheSys)},
		{name: "MSpanInuse", value: float64(rtm.MSpanInuse)},
		{name: "MSpanSys", value: float64(rtm.MSpanSys)},
		{name: "Mallocs", value: float64(rtm.Mallocs)},
		{name: "NextGC", value: float64(rtm.NextGC)},
		{name: "NumForcedGC", value: float64(rtm.NumForcedGC)},
		{name: "NumGC", value: float64(rtm.NumGC)},
		{name: "OtherSys", value: float64(rtm.OtherSys)},
		{name: "PauseTotalNs", value: float64(rtm.PauseTotalNs)},
		{name: "StackInuse", value: float64(rtm.StackInuse)},
		{name: "StackSys", value: float64(rtm.StackSys)},
		{name: "Sys", value: float64(rtm.Sys)},
	}
}

//
//func schedule(f func(), interval time.Duration) *time.Ticker {
//	ticker := time.NewTicker(interval)
//	go func() {
//		for range ticker.C {
//			f()
//		}
//	}()
//	return ticker
//}

func GetRuntimeMetricsAsync(interval time.Duration) *time.Ticker {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			GetRuntimeMetrics()
		}
	}()
	return ticker
}

func main() {
	go GetRuntimeMetricsAsync(pollRuntimeMetricsInterval)
	go SendDataAsync()

	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-cancelSignal

	os.Exit(1)
}
