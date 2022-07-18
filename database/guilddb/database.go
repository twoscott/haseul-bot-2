package guilddb

import (
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/jmoiron/sqlx"
	"github.com/twoscott/haseul-bot-2/router"
)

// DB wraps an sqlx database instance with helper methods for guilds querying.
type DB struct {
	*sqlx.DB
}

// New returns a new instance of a guilds database.
func New(dbConn *sqlx.DB) *DB {
	db := &DB{dbConn}
	db.createTables()
	return db
}

func (db *DB) createTables() {
	db.MustExec(createGuildConfigsTableQuery)
}

func (db *DB) Init(rt *router.Router) {
	rt.AddStartupListener(db.onStartup)
}

func (db *DB) onStartup(rt *router.Router, _ *gateway.ReadyEvent) {
	guilds, _ := rt.State.Guilds()
	for _, guild := range guilds {
		db.Add(guild.ID)
	}
}
