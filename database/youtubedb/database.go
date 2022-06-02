package youtubedb

import "github.com/jmoiron/sqlx"

// DB wraps an sqlx database instance with helper methods for Twitter querying.
type DB struct {
	*sqlx.DB
}

// New returns a new instance of a Twitter database.
func New(dbConn *sqlx.DB) *DB {
	db := &DB{dbConn}
	db.createTables()
	return db
}

func (db *DB) createTables() {
	db.MustExec(createYouTubeHistoryTable)
	db.MustExec(createHistoryToggleTableQuery)
}
