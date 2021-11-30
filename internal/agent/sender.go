package agent

import (
	"fmt"
	"log"
	"net/http"
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
	url := fmt.Sprintf("%s%s:%s/update/%s/%s/%d",
		DefaultProtocol,
		DefaultServer,
		DefaultPort,
		"gauge",
		g.name,
		int(g.value))
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

func (c Counter) SendCounter(client *http.Client) (bool, error) {
	url := fmt.Sprintf("%s%s:%s/update/%s/%s/%d",
		DefaultProtocol,
		DefaultServer,
		DefaultPort,
		"conter",
		c.name,
		int(c.value))
	log.Println("SendCounter " + url)
	resp, err := client.Post(url, ContentType, nil)
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
