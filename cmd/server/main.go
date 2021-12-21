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
	addressFlag, storeFileFlag string
	storeIntervalFlag          time.Duration
	restoreFlag                bool
)

const (
	defaultAddress       = "localhost:8080"
	defaultStoreInterval = 300 * time.Second
	defaultStoreFile     = "/tmp/devops-metrics-db.json"
	defaultRestore       = true
)

func getConfig() types.ServerConfig {
	var c types.ServerConfig
	err := env.Parse(&c)
	if err != nil {
		log.Println("Can't read env config")
		log.Println(err)
	}

	flag.StringVar(&addressFlag, "a", defaultAddress, "Server Address")
	flag.StringVar(&storeFileFlag, "f", defaultStoreFile, "File Path")
	flag.DurationVar(&storeIntervalFlag, "i", defaultStoreInterval, "Store Interval")
	flag.BoolVar(&restoreFlag, "r", defaultRestore, "Restore After Start")
	flag.Parse()

	//rewrite if env values is not empty
	_, isSet := os.LookupEnv("ADDRESS")
	if !isSet {
		c.ServerAddress = addressFlag
	}

	_, isSet = os.LookupEnv("STORE_FILE")
	if !isSet {
		c.FileStoragePath = storeFileFlag
	}

	_, isSet = os.LookupEnv("STORE_INTERVAL")
	if !isSet {
		c.StoreInterval = storeIntervalFlag
	}
	_, isSet = os.LookupEnv("RESTORE")
	if !isSet {
		c.Restore = restoreFlag
	}

	return c
}

func StartServer(c types.ServerConfig) {

	if c.Restore && c.FileStoragePath != "" {
		server.LoadMetrics(c)
	}

	if c.StoreInterval > 0 && c.FileStoragePath != "" {
		go server.SaveMetrics(c)
	}

	svr := &http.Server{
		Addr:    addressFlag,
		Handler: server.Router(),
	}
	svr.SetKeepAlivesEnabled(false)
	log.Printf("listening on port %s ", addressFlag)
	log.Fatal(svr.ListenAndServe())

}

func main() {

	types.SConfig = getConfig()
	types.SConfig.LogConfig()

	go StartServer(types.SConfig)

	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-cancelSignal
}
