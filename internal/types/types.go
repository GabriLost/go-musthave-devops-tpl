package types

import (
	"log"
	"time"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type InternalStorage struct {
	CounterMetrics map[string]int64
	GaugeMetrics   map[string]float64
}

type (
	AgentConfig struct {
		Address        string        `env:"ADDRESS"`
		PollInterval   time.Duration `env:"POLL_INTERVAL"`
		ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	}
	ServerConfig struct {
		ServerAddress   string        `env:"ADDRESS"`
		FileStoragePath string        `env:"STORE_FILE"`
		StoreInterval   time.Duration `env:"STORE_INTERVAL"`
		Restore         bool          `env:"RESTORE"`
	}
)

var (
	SenderConfig = AgentConfig{
		Address:        "localhost:8080",
		PollInterval:   2,
		ReportInterval: 5,
	}
	SConfig = ServerConfig{}
)

func (c AgentConfig) LogConfig() {
	log.Printf(`agent address="%s", poll interval="%s" report interval="%s"`,
		c.Address, c.PollInterval, c.ReportInterval)
}

func (c ServerConfig) LogConfig() {
	log.Printf(`server address="%s", file path="%s", store interval="%s", is restore="%t"`,
		c.ServerAddress, c.FileStoragePath, c.StoreInterval, c.Restore)
}
