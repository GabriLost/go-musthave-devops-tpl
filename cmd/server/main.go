package main

import (
	"github.com/GabriLost/go-musthave-devops-tpl/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultServer = "localhost"
	defaultPort   = "8080"
)

func StartServer() {
	svr := &http.Server{
		Addr:    defaultServer + ":" + defaultPort,
		Handler: server.Router(),
	}
	svr.SetKeepAlivesEnabled(false)
	log.Printf("listening on port " + defaultPort)
	log.Fatal(svr.ListenAndServe())
}

func main() {
	go StartServer()

	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-cancelSignal
}
