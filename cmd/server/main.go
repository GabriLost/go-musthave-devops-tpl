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
	addressFlag, storeFileFlag, keyFlag, dbDNSFlag string
	storeIntervalFlag                              time.Duration
	restoreFlag                                    bool
)

const (
	defaultAddress       = "localhost:8080"
	defaultStoreInterval = 300 * time.Second
	defaultStoreFile     = "/tmp/devops-metrics-db.json"
	defaultRestore       = true
)

func getConfig() (types.ServerConfig, error) {
	var c types.ServerConfig
	err := env.Parse(&c)
	if err != nil {
		return c, err
	}

	flag.StringVar(&addressFlag, "a", defaultAddress, "Server Address")
	flag.StringVar(&storeFileFlag, "f", defaultStoreFile, "File Path")
	flag.DurationVar(&storeIntervalFlag, "i", defaultStoreInterval, "Store Interval")
	flag.BoolVar(&restoreFlag, "r", defaultRestore, "Restore After Start")
	flag.StringVar(&keyFlag, "k", "", "Secret Key")
	flag.StringVar(&dbDNSFlag, "d", "", "Database DNS")

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

	_, isSet = os.LookupEnv("KEY")
	if !isSet {
		c.Key = keyFlag
	}

	_, isSet = os.LookupEnv("DATABASE_DSN")
	if !isSet {
		c.DatabaseDSN = dbDNSFlag
	}

	return c, nil
}

func StartServer(c types.ServerConfig) {

	err := server.LoadIndexHTML()
	if err != nil {
		log.Println("index page not loaded")
	}

	if err := server.ConnectDB(); err != nil {
		log.Printf("failed to connect db: %v", err)
	}

	if c.Restore {
		// db storage has priority
		if c.DatabaseDSN != "" {
			if err := server.LoadStatsDB(); err != nil {
				log.Print(err)
			}
		} else if c.FileStoragePath != "" {
			if err := server.LoadMetrics(c); err != nil {
				log.Print(err)
			}
		}
	}
	//save into file, if DatabaseDSN is empty
	if c.DatabaseDSN == "" && c.StoreInterval > 0 && c.FileStoragePath != "" {
		go server.SaveMetricsIntoFileBySchedule(c)
	}

	svr := &http.Server{
		Addr:    c.ServerAddress,
		Handler: server.Router(),
	}
	svr.SetKeepAlivesEnabled(false)
	log.Printf("listening on port %s ", c.ServerAddress)
	log.Fatal(svr.ListenAndServe())

}

func main() {

	cfg, err := getConfig()
	if err != nil {
		log.Fatal()
	}
	types.SConfig = cfg
	types.SConfig.LogConfig()

	go StartServer(types.SConfig)

	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-cancelSignal
}
