package server

import (
	"database/sql"
	"github.com/GabriLost/go-musthave-devops-tpl/internal/types"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose"
	"log"
)

var DB *sql.DB

func ConnectDB() (err error) {
	DB, err = sql.Open("pgx", types.SConfig.DatabaseDSN)
	if err != nil {
		log.Printf("unable to connect to database: %s", err)
		return err
	}
	log.Printf("Database connection was created")

	log.Printf("Start migrating database \n")

	err = goose.Up(DB, "internal/migrations")
	if err != nil {
		return err
	}

	return nil

}
func SaveGaugeDB(name string, value float64) error {
	_, err := DB.Exec("INSERT INTO gauges (name, value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE set value = $2", name, value)
	if err != nil {
		return err
	}
	return nil
}

func SaveCounterDB(name string, delta int64) error {
	_, err := DB.Exec("INSERT INTO counters (name, value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE SET value = $2", name, delta)
	if err != nil {
		return err
	}
	return nil
}

func LoadStatsDB() error {
	var name string
	var gauge float64
	var counter int64

	gRows, err := DB.Query("SELECT name, value FROM gauges")
	if err != nil {
		return err
	}
	defer gRows.Close()
	for gRows.Next() {
		if err = gRows.Scan(&name, &gauge); err != nil {
			log.Print(err)
			return err
		}
		MetricGauges[name] = gauge
	}
	if err = gRows.Err(); err != nil {
		return err
	}

	cRows, err := DB.Query("SELECT name, value FROM counters")
	if err != nil {
		return err
	}
	defer cRows.Close()
	for cRows.Next() {
		if err = cRows.Scan(&name, &counter); err != nil {
			log.Print(err)
			return err
		}
		MetricCounters[name] = counter
	}
	if err = cRows.Err(); err != nil {
		return err
	}

	return nil
}
