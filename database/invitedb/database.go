package invitedb

import (
	"github.com/jmoiron/sqlx"
)

// DB wraps an sqlx database instance with helper methods for invite querying.
type DB struct {
	*sqlx.DB
}

// New returns a new instance of an invites database.
func New(dbConn *sqlx.DB) *DB {
	db := &DB{dbConn}
	db.createTables()
	return db
}

func (db *DB) createTables() {
	db.MustExec(createGuildInvitesTableQuery)
}
