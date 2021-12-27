package server

import (
	"encoding/json"
	"github.com/GabriLost/go-musthave-devops-tpl/internal/types"
	"log"
	"os"
	"time"
)

var (
	MetricCounters = make(map[string]int64)
	MetricGauges   = make(map[string]float64)
)

func LoadMetrics(c types.ServerConfig) {
	log.Printf("Loading metrics from file %s", c.FileStoragePath)

	flag := os.O_RDONLY

	f, err := os.OpenFile(c.FileStoragePath, flag, 0)
	if err != nil {
		log.Print("Can't open file for loading metrics: ", err)
		return
	}
	defer f.Close()

	var internalStorage types.InternalStorage

	if err := json.NewDecoder(f).Decode(&internalStorage); err != nil {
		log.Fatal("Can't decode metrics: ", err)
	}

	MetricGauges = internalStorage.GaugeMetrics
	MetricCounters = internalStorage.CounterMetrics
	log.Printf("Metrics successfully loaded from file %s", c.FileStoragePath)
}

func SaveMetrics(c types.ServerConfig) {
	ticker := time.NewTicker(c.StoreInterval)
	for {
		<-ticker.C
		log.Printf("Dumping metrics to file %s", c.FileStoragePath)
		SaveMetricsImpl(c)
	}
}

func SaveMetricsImpl(c types.ServerConfig) {
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	f, err := os.OpenFile(c.FileStoragePath, flag, 0644)
	if err != nil {
		log.Fatal("Can't open file for dumping: ", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)

	internalStorage := types.InternalStorage{
		GaugeMetrics:   MetricGauges,
		CounterMetrics: MetricCounters,
	}

	if err := encoder.Encode(internalStorage); err != nil {
		log.Fatal("Can't encode server's metrics: ", err)
	}
}
