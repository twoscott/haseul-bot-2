package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/twoscott/haseul-bot-2/config"
)

func mustGetConnection() *sqlx.DB {
	cfg := config.GetInstance()
	dbName := cfg.PostgreSQL.Database
	user := cfg.PostgreSQL.Username
	password := cfg.PostgreSQL.Password

	connStr := fmt.Sprintf(
		"dbname=%s user=%s password=%s sslmode=disable",
		dbName, user, password,
	)
	dbConn, err := sqlx.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = dbConn.Ping()
	if err != nil {
		panic(err)
	}

	return dbConn
}
