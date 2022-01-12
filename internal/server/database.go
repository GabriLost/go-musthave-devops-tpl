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

	err = goose.Up(DB, ".")
	if err != nil {
		return err
	}

	return nil

}
