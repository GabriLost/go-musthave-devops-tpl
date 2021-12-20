package types

import (
	"time"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

var (
	MetricCounters = make(map[string]int64)
	MetricGauges   = make(map[string]float64)
)

type (
	AgentConfig struct {
		Address        string        `env:"ADDRESS" envDefault:"localhost1"`
		PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"1s"`
		ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"5s"`
	}
	ServerConfig struct {
		ServerAddress   string        `env:"SERVER_ADDRESS" envDefault:"localhost"`
		FileStoragePath string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
		StoreInterval   time.Duration `env:"STORE_INTERVAL" envDefault:"300"`
		Restore         bool          `env:"RESTORE" envDefault:"true"`
	}
)

var (
	// SenderConfig config for sender service
	SenderConfig = AgentConfig{
		Address:        "localhost2",
		PollInterval:   5,
		ReportInterval: 15,
	}

	SConfig = ServerConfig{}
)
