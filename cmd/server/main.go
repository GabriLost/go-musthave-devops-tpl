package main

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"text/template"
)

const (
	defaultServer = "localhost"
	defaultPort   = "8080"
	contentType   = "application/json"
)

var metrics = make(map[string]string)

var templateDataMap = make(map[string]interface{})

const HTMLPage = `
<div>
	<h1>Metric monitor</h1>
	<table>
		<thead>
			<tr>
				<th>Name</th>
				<th>Value </th>
			</tr>
		</thead>
		<tbody>
			{{ range $key, $val :=  .metrics }}
				<tr>
					<td>{{$key}}</td>
					<td>{{$val}}</td>
				</tr>
			{{end}}
		</tbody>
	</table>
</div>
`

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	templateDataMap["metrics"] = metrics
	tmpl := template.Must(template.New("").Parse(HTMLPage))
	err := tmpl.Execute(w, templateDataMap)
	if err != nil {
		return
	}
}

func GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI

	if len(strings.Split(uri, "/")) != 5 {
		http.Error(w, "URI too long or too short", http.StatusBadRequest)
		return
	}
	metric := strings.Split(uri, "/")[2]
	name := strings.Split(uri, "/")[3]
	value := strings.Split(uri, "/")[4]
	_, err1 := strconv.ParseFloat(value, 64)
	_, err2 := strconv.ParseInt(value, 10, 64)
	if err1 != nil && err2 != nil {
		http.Error(w, "parsing error", http.StatusBadRequest)
		return
	} else if metric != "gauge" && metric != "counter" {
		http.Error(w, "No such type of metric", http.StatusBadRequest)
		return
	} else {
		metrics[name] = value
		w.Header().Add("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
	}
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "NotImplemented", http.StatusNotImplemented)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "StatusNotFound", http.StatusNotFound)
}

func StartServer() {

	router := chi.NewRouter()
	router.Get("/", GetAllHandler)
	router.Post("/*", NotFound)
	router.Post("/update/*", NotImplemented)
	router.Post("/update/gauge/{name}/{value}", GetMetricHandler)
	router.Post("/update/counter/{name}/{value}", GetMetricHandler)

	server := &http.Server{
		Addr:    defaultServer + ":" + defaultPort,
		Handler: router,
	}
	server.SetKeepAlivesEnabled(false)
	log.Printf("listening on port " + defaultPort)
	log.Fatal(server.ListenAndServe())
}

func main() {
	go StartServer()

	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-cancelSignal
}
