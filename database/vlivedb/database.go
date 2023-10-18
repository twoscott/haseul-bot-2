package vlivedb

import "github.com/jmoiron/sqlx"

// DB wraps an sqlx database instance with helper methods for VLIVE querying.
type DB struct {
	*sqlx.DB
}

// New returns a new instance of a VLIVE database.
func New(dbConn *sqlx.DB) *DB {
	db := &DB{dbConn}
	db.createTables()
	return db
}

func (db *DB) createTables() {
	db.MustExec(createVLIVEBoardsTableQuery)
	db.MustExec(createVLIVEFeedsTableQuery)
	db.MustExec(createVLIVEMentionsTableQuery)
	db.MustExec(createVLIVERetriesTableQuery)
}
