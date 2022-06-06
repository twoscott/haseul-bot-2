package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/twoscott/haseul-bot-2/config"
)

func mustGetConnection() *sqlx.DB {
	cfg := config.GetInstance()
	dbName := cfg.PostgreSQL.Database
	user := cfg.PostgreSQL.Username
	password := cfg.PostgreSQL.Password
	if dbName == "" || user == "" || password == "" {
		log.Fatalln("Invalid database config variables provided")
	}

	connStr := fmt.Sprintf(
		"dbname=%s user=%s password=%s sslmode=disable",
		dbName, user, password,
	)
	dbConn, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s\n", err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %s\n", err)
	}

	return dbConn
}
