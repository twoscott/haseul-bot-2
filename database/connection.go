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

	connStr := fmt.Sprintf(
		"host=%s "+
			"port=%s "+
			"dbname=%s "+
			"user=%s "+
			"password=%s "+
			"sslmode=disable",
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.Database,
		cfg.PostgreSQL.Username,
		cfg.PostgreSQL.Password,
	)
	dbConn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s\n", err)
	}

	return dbConn
}
