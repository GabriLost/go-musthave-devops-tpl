package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/GabriLost/go-musthave-devops-tpl/internal/types"
	"log"
	"net/http"
	"strings"
)

type Gauge struct {
	name  string
	value float64
}

type Counter struct {
	name  string
	value int64
}

func (g Gauge) SendGauge(client *http.Client) (bool, error) {
	url := fmt.Sprintf("%s%s/update/",
		DefaultProtocol,
		types.SenderConfig.Address)

	var m = types.Metrics{
		ID:    g.name,
		MType: "gauge",
		Value: &g.value,
		Hash:  "",
	}
	m.AddHashWithKey(types.SenderConfig.Key)
	b, _ := json.Marshal(m)

	log.Printf("SendGauge %s %s", string(b), url)
	body := strings.NewReader(string(b))
	resp, err := client.Post(url, "application/json", body)
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

func (c Counter) SendCounter(client *http.Client) (bool, error) {
	url := fmt.Sprintf("%s%s/update/",
		DefaultProtocol,
		types.SenderConfig.Address)

	var m = types.Metrics{
		ID:    c.name,
		MType: "counter",
		Delta: &c.value,
		Hash:  "",
	}
	m.AddHashWithKey(types.SenderConfig.Key)
	b, _ := json.Marshal(m)

	log.Printf("SendCounter %s %s", string(b), url)

	body := strings.NewReader(string(b))
	resp, err := client.Post(url, "application/json", body)
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

func appendBatch(initial []types.Metrics, name string, data interface{}) []types.Metrics {
	if initial == nil {
		log.Print("addBatch: trying to add to nil slice")
	}
	switch v := data.(type) {
	case int64:
		delta := v
		m := types.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &delta,
		}
		m.AddHashWithKey(types.SenderConfig.Key)
		return append(initial, m)
	case float64:
		value := v
		m := types.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		}
		m.AddHashWithKey(types.SenderConfig.Key)
		return append(initial, m)
	default:
		return initial
	}
}

func sendMetricsBatch(metrics []types.Metrics) {

	log.Println("Total metrics in Batch is", len(metrics))

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(metrics); err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("%s%s/updates/",
		DefaultProtocol,
		types.SenderConfig.Address)
	resp, err := http.Post(url, "application/json", &body)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Sending %s, http status %d", url, resp.StatusCode)
	}
}

func SendMetrics() {
	Metrics = StoreRandomMetrics(Metrics)

	Metrics = append(Metrics, UtilizationData.TotalMemory)
	Metrics = append(Metrics, UtilizationData.FreeMemory)
	Metrics = append(Metrics, UtilizationData.CPUutilizations...)

	log.Println("Total metrics is", len(Metrics))
	client := http.Client{Timeout: DefaultTimeout}
	if types.SenderConfig.UseBatch {
		metrics := make([]types.Metrics, 0, len(Metrics))
		for _, m := range Metrics {
			metrics = appendBatch(metrics, m.name, m.value)
		}
		metricCounter := Counter{name: "PollCount", value: PollCount}
		metrics = appendBatch(metrics, metricCounter.name, metricCounter.value)
		sendMetricsBatch(metrics)
	} else {
		for _, i := range Metrics {
			_, err := i.SendGauge(&client)
			if err != nil {
				log.Println("can't send Gauge " + err.Error())
				return
			}
		}
		metricCounter := Counter{name: "PollCount", value: PollCount}
		_, err := metricCounter.SendCounter(&client)
		if err != nil {
			log.Println("can't send Counter " + err.Error())
			return
		}
	}

	log.Println("Reset poll counter to zero")
	PollCount = 0

	client.CloseIdleConnections()
}
