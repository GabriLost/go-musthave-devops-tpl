package types

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
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
		Key            string        `env:"KEY"`
	}
	ServerConfig struct {
		ServerAddress   string        `env:"ADDRESS"`
		FileStoragePath string        `env:"STORE_FILE"`
		StoreInterval   time.Duration `env:"STORE_INTERVAL"`
		Restore         bool          `env:"RESTORE"`
		Key             string        `env:"KEY"`
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
	log.Printf(`agent address="%s", poll interval="%s", report interval="%s", key="%s"`,
		c.Address, c.PollInterval, c.ReportInterval, c.Key)
}

func (c ServerConfig) LogConfig() {
	log.Printf(`server address="%s", file path="%s", store interval="%s", is restore="%t", key="%s`,
		c.ServerAddress, c.FileStoragePath, c.StoreInterval, c.Restore, c.Key)
}

func (m *Metrics) AddHashWithKey(key string) error {
	if key == "" {
		return nil
	}

	h, err := m.computeHashWithKey(key)
	if err != nil {
		return err
	}

	m.Hash = hex.EncodeToString(*h)

	return nil
}

func (m Metrics) CheckHashWithKey(key string) error {
	if key == "" {
		return nil
	}

	h, err := m.computeHashWithKey(key)
	if err != nil {
		return err
	}

	hashStr := hex.EncodeToString(*h)
	if m.Hash != hashStr {
		return errors.New("incorrect hash")
	}

	return nil
}

func (m Metrics) computeHashWithKey(key string) (*[]byte, error) {
	if key == "" {
		return nil, fmt.Errorf("no key")
	}
	if m.ID == "" {
		return nil, fmt.Errorf("empty ID field")
	}
	data := ""
	switch m.MType {
	case "gauge":
		{
			data = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
		}
	case "counter":
		{
			data = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
		}
	default:
		return nil, errors.New("no such type")
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	hash := h.Sum(nil)

	return &hash, nil
}
