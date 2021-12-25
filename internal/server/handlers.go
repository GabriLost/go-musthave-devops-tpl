package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GabriLost/go-musthave-devops-tpl/internal/types"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

var HTMLTemplate *template.Template

func AllMetricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	data := make(map[string]interface{})
	data[MetricTypeGauge] = MetricGauges
	data[MetricTypeCounter] = MetricCounters
	err := HTMLTemplate.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}

func LoadIndexHTML() error {
	bytes, err := os.ReadFile("internal/server/" + HTMLFile)
	if err != nil {
		bytes, err = os.ReadFile(HTMLFile)
		if err != nil {
			return err
		}
	}

	HTMLTemplate, err = template.New("").Parse(string(bytes))
	if err != nil {
		return err
	}
	return nil
}

func UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {

	metric := chi.URLParam(r, "typ")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	switch metric {
	case MetricTypeGauge:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, "Wrong gauge value", http.StatusBadRequest)
			return
		}
		MetricGauges[name] = val
	case MetricTypeCounter:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, "Wrong counter value", http.StatusBadRequest)
			return
		}
		MetricCounters[name] += val
	default:
		http.Error(w, "No such type of metric", http.StatusNotImplemented)
		return
	}
	log.Printf("got metric %s", name)

	w.Header().Add("Content-Type", contentTypeAppJSON)
	w.WriteHeader(http.StatusOK)
}

func ValueMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "typ")
	metricName := chi.URLParam(r, "name")
	switch metricType {
	case MetricTypeGauge:
		if val, found := MetricGauges[metricName]; found {
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", contentTypeAppJSON)
			_, err := w.Write([]byte(fmt.Sprint(val)))
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			http.Error(w, "There is no metric you requested", http.StatusNotFound)
		}
	case MetricTypeCounter:
		if val, found := MetricCounters[metricName]; found {
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", contentTypeAppJSON)
			_, err := w.Write([]byte(fmt.Sprint(val)))
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			http.Error(w, "There is no metric you requested", http.StatusNotFound)
		}
	default:
		http.Error(w, "There is no metric you requested", http.StatusBadRequest)
	}
}

func NotImplementedHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Method is not implemented yet", http.StatusNotImplemented)
}

func NotFoundHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad Request"))
}

func JSONUpdateMetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != contentTypeAppJSON {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"Status":"Bad Request"}`))
		if err != nil {
			log.Println("Wrong content type")
			return
		}
		return
	}

	var m types.Metrics
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"Status":"Bad Request"}`))
		if err != nil {
			log.Println("Decode problem")
			return
		}
		return
	}

	if err := m.CheckHashWithKey(types.SConfig.Key); err != nil {
		ResponseErrorJSON(w, http.StatusBadRequest, "Incorrect Hash")
		return
	}

	//todo validate
	err := saveMetrics(m)
	if err != nil {
		ResponseErrorJSON(w, http.StatusNotImplemented, "No such type of metric")
		return
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		ResponseErrorJSON(w, http.StatusInternalServerError, "can't encode metrics")
		return
	}
	w.Header().Add("Content-Type", contentTypeAppJSON)
}

func JSONValueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		ResponseErrorJSON(w, http.StatusBadRequest, "Header type is not \"application/json\"")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		ResponseErrorJSON(w, http.StatusInternalServerError, "can't read body")
		return
	}

	var m types.Metrics

	err = json.Unmarshal(body, &m)
	if err != nil {
		ResponseErrorJSON(w, http.StatusBadRequest, "can't Unmarshal body")
		return
	}

	if m.ID == "" {
		ResponseErrorJSON(w, http.StatusBadRequest, "ID(metricName) is null")
		return
	}

	switch m.MType {
	case MetricTypeGauge:
		val, ok := MetricGauges[m.ID]
		if !ok {
			ResponseErrorJSON(w, http.StatusNotFound, "MetricTypeGauge "+m.ID)
			return
		}
		m.Value = &val
	case MetricTypeCounter:
		val, ok := MetricCounters[m.ID]
		if !ok {
			ResponseErrorJSON(w, http.StatusNotFound, "MetricTypeCounter "+m.ID)
			return
		}
		m.Delta = &val
	default:
		ResponseErrorJSON(w, http.StatusNotFound, "metric type not found "+m.MType)
		return
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		ResponseErrorJSON(w, http.StatusInternalServerError, "can't encode metrics")
		return
	}

}

func ResponseErrorJSON(w http.ResponseWriter, statusCode int, message string) {
	log.Printf("ResponseErrorJSON with status code %d, %s", statusCode, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	statusStr := http.StatusText(statusCode)
	if statusStr == "" {
		statusStr = "UNKNOWN"
	}
	type Response struct {
		Status string `json:"Status"`
	}

	r := Response{Status: statusStr}
	err := json.NewEncoder(w).Encode(&r)
	if err != nil {
		return
	}
}

func saveMetrics(m types.Metrics) error {
	err := m.AddHashWithKey(types.SConfig.Key)
	if err != nil {
		return err
	}
	switch m.MType {
	case MetricTypeGauge:
		log.Printf("saving metric %s %s %f\n", m.ID, m.MType, *m.Value)
		MetricGauges[m.ID] = *m.Value
	case MetricTypeCounter:
		log.Printf("saving metric %s %s %d\n", m.ID, m.MType, *m.Delta)
		MetricCounters[m.ID] += *m.Delta
	default:
		return errors.New("no such type of metric")
	}
	return nil

}
