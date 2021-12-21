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

func GetAllHandler(w http.ResponseWriter, _ *http.Request) {
	//todo спросить почему так?
	indexPage, err := os.ReadFile("internal/server/index.html")
	if err != nil {
		indexPage, err = os.ReadFile("index.html")
		if err != nil {
			log.Println(err)
			return
		}
	}
	indexTemplate := template.Must(template.New("").Parse(string(indexPage)))
	tmp := make(map[string]interface{})
	tmp[MetricTypeGauge] = MetricGauges
	tmp[MetricTypeCounter] = MetricCounters
	err = indexTemplate.Execute(w, tmp)
	if err != nil {
		log.Println(err)
		return
	}
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

	//todo validate
	err := saveMetrics(m)
	if err != nil {
		w.Header().Set("Content-Type", contentTypeAppJSON)
		log.Println(err)
		http.Error(w, "No such type of metric", http.StatusNotImplemented)
		return
	}
	w.Header().Add("Content-Type", contentTypeAppJSON)
	w.WriteHeader(http.StatusOK)

}

func JSONValueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		ResponseErrorJSON(w, http.StatusBadRequest)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var m types.Metrics

	err := json.Unmarshal(body, &m)
	if err != nil {
		log.Println(err)
		ResponseErrorJSON(w, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &m)
	if err != nil {
		ResponseErrorJSON(w, http.StatusInternalServerError)
		return
	}
	// if ID is null, then bad request
	if m.ID == "" {
		ResponseErrorJSON(w, http.StatusBadRequest)
		return
	}

	switch m.MType {
	case MetricTypeGauge:
		val, ok := MetricGauges[m.ID]
		if !ok {
			ResponseErrorJSON(w, http.StatusNotFound)
			return
		}
		m.Value = &val
	case MetricTypeCounter:
		val, ok := MetricCounters[m.ID]
		if !ok {
			ResponseErrorJSON(w, http.StatusNotFound)
			return
		}
		m.Delta = &val
	default:
		ResponseErrorJSON(w, http.StatusNotFound)
		return
	}
	// encode
	ret, err := json.Marshal(m)
	if err != nil {
		ResponseErrorJSON(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

func ResponseErrorJSON(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	statusStr := http.StatusText(statusCode)
	if statusStr == "" {
		statusStr = "UNKNOWN"
	}
	_, err := w.Write([]byte(`{"Status":"` + statusStr + `"}`))
	if err != nil {
		return
	}
}

func saveMetrics(m types.Metrics) error {
	//todo mutex
	log.Printf("saving metric %s %s %d %d\n", m.ID, m.MType, m.Delta, m.Value)
	switch m.MType {
	case MetricTypeGauge:
		MetricGauges[m.ID] = *m.Value
	case MetricTypeCounter:
		MetricCounters[m.ID] += *m.Delta
	default:
		return errors.New("no such type of metric")
	}
	return nil

}
