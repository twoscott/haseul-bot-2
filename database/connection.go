package database

import (
	"fmt"

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

	return sqlx.MustConnect("postgres", connStr)
}
