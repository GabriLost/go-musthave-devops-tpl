package agent

import (
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

func SendMetrics() {
	Metrics = StoreRandomMetrics(Metrics)
	log.Println("Total metrics is ", len(Metrics))
	client := http.Client{Timeout: DefaultTimeout}
	for _, i := range Metrics {
		_, err := i.SendGauge(&client)
		if err != nil {
			log.Println("can't send Gauge " + err.Error())
			return
		}
	}
	metricCounter := Counter{name: "PollCount", value: PollCount}
	log.Println("Reset poll counter to zero")
	PollCount = 0
	_, err := metricCounter.SendCounter(&client)
	if err != nil {
		log.Println("can't send Counter " + err.Error())
		return
	}
	client.CloseIdleConnections()
}
