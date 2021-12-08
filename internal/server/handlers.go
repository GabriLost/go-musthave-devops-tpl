package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GabriLost/go-musthave-devops-tpl/internal/types"
	"github.com/go-chi/chi/v5"
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
		log.Println(err)
		os.Exit(1)
	}
	indexTemplate := template.Must(template.New("").Parse(string(indexPage)))
	tmp := make(map[string]interface{})
	tmp[MetricTypeGauge] = types.MetricGauges
	tmp[MetricTypeCounter] = types.MetricCounters
	err = indexTemplate.Execute(w, tmp)
	if err != nil {
		log.Println(err)
		return
	}
}

func PostMetricHandler(w http.ResponseWriter, r *http.Request) {

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
		types.MetricGauges[name] = val
	case MetricTypeCounter:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, "Wrong counter value", http.StatusBadRequest)
			return
		}
		types.MetricCounters[name] += val
	default:
		http.Error(w, "No such type of metric", http.StatusNotImplemented)
		return
	}
	log.Printf("got metric %s", name)

	w.Header().Add("Content-Type", contentTypeAppJson)
	w.WriteHeader(http.StatusOK)
}

func GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "typ")
	metricName := chi.URLParam(r, "name")
	switch metricType {
	case MetricTypeGauge:
		if val, found := types.MetricGauges[metricName]; found {
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", contentTypeAppJson)
			_, err := w.Write([]byte(fmt.Sprint(val)))
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			http.Error(w, "There is no metric you requested", http.StatusNotFound)
		}
	case MetricTypeCounter:
		if val, found := types.MetricCounters[metricName]; found {
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", contentTypeAppJson)
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

func PostJsonMetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != contentTypeAppJson {
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
		w.Header().Set("Content-Type", contentTypeAppJson)
		log.Println(err)
		http.Error(w, "No such type of metric", http.StatusNotImplemented)
		return
	}
	w.Header().Add("Content-Type", contentTypeAppJson)
	w.WriteHeader(http.StatusOK)

}

func saveMetrics(m types.Metrics) error {
	//todo mutex
	log.Printf("%s %s %d %d\n", m.ID, m.MType, m.Delta, m.Value)
	switch m.MType {
	case MetricTypeGauge:
		types.MetricGauges[m.ID] = *m.Value
	case MetricTypeCounter:
		types.MetricCounters[m.ID] += *m.Delta
	default:
		return errors.New("No such type of metric ")
	}
	return nil

}
