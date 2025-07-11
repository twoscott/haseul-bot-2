package marriagedb

import "github.com/jmoiron/sqlx"

// DB wraps an sqlx database instance with helper methods for
// Command querying.
type DB struct {
	*sqlx.DB
}

// New returns a new instance of a Marriages database.
func New(dbConn *sqlx.DB) *DB {
	db := &DB{dbConn}
	db.createTables()
	return db
}

func (db *DB) createTables() {
	db.MustExec(createMarriagesTableQuery)
}
