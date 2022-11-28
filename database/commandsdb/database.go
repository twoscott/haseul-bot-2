package commandsdb

import "github.com/jmoiron/sqlx"

// DB wraps an sqlx database instance with helper methods for
// Notification querying.
type DB struct {
	*sqlx.DB
}

// New returns a new instance of a Notifications database.
func New(dbConn *sqlx.DB) *DB {
	db := &DB{dbConn}
	db.createTables()
	return db
}

func (db *DB) createTables() {
	db.MustExec(createCommandsTableQuery)
}
