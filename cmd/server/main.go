package main

import (
	"flag"
	"github.com/GabriLost/go-musthave-devops-tpl/internal/server"
	"github.com/GabriLost/go-musthave-devops-tpl/internal/types"
	"github.com/caarlos0/env/v6"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	serverAddress, storeFile string
	storeInterval            int
	restore                  bool
)

func StartServer() {

	var cfg types.ServerConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Println("Can't read env config")
		log.Println(err)
	}

	//rewrite if flags is not empty
	flag.StringVar(&serverAddress, "a", "", "Server Address")
	flag.StringVar(&storeFile, "f", "", "File Path")
	flag.IntVar(&storeInterval, "i", -1, "Store Interval")
	flag.BoolVar(&restore, "r", false, "Restore After Start")
	flag.Parse()

	if serverAddress != "" {
		types.SConfig.ServerAddress = serverAddress
	}
	if storeFile != "" {
		types.SConfig.FileStoragePath = storeFile
	}
	if storeInterval != -1 {
		types.SConfig.StoreInterval = time.Second * time.Duration(storeInterval)
	}

	if !restore {
		types.SConfig.Restore = restore
	}

	svr := &http.Server{
		Addr:    serverAddress + ":" + "8080",
		Handler: server.Router(),
	}
	svr.SetKeepAlivesEnabled(false)
	log.Printf("listening on port " + "8080")
	log.Fatal(svr.ListenAndServe())
}

func main() {
	go StartServer()

	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-cancelSignal
}
