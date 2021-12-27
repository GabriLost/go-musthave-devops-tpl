package agent

import (
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
	b := fmt.Sprintf(`{"id":"%s", "type":"%s", "value": %d}`,
		g.name,
		"gauge",
		int(g.value))
	log.Printf("SendGauge %s %s", b, url)
	body := strings.NewReader(b)
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
	b := fmt.Sprintf(`{"id":"%s", "type":"%s", "delta": %d}`,
		c.name,
		"counter",
		c.value)
	log.Printf("SendCounter %s %s", b, url)
	body := strings.NewReader(b)
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
