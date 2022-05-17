package twitterdb

const (
	createTwitterRetriesTableQuery = `
		CREATE TABLE IF NOT EXISTS TwitterRetries(
			tweetID INT8 NOT NULL,
			PRIMARY KEY(tweetID)
		)`

	addRetryQuery = `
		INSERT INTO TwitterRetries VALUES($1) ON CONFLICT DO NOTHING`
	getAllRetriesQuery = `
		SELECT tweetID FROM TwitterRetries`
	removeRetryQuery = `
		DELETE FROM TwitterRetries WHERE tweetID = $1`
)

// AddRetry returns adds a tweet to retry sending to the database.
func (db *DB) AddRetry(tweetID int64) (bool, error) {
	res, err := db.Exec(addRetryQuery, tweetID)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// GetAllUsers returns all Tweets that previously failed to send from
// the database.
func (db *DB) GetAllRetries() ([]int64, error) {
	var tweetIDs []int64
	err := db.Select(&tweetIDs, getAllRetriesQuery)

	return tweetIDs, err
}

// RemoveRetry removes a retry tweet from the database.
func (db *DB) RemoveRetry(tweetID int64) (bool, error) {
	res, err := db.Exec(removeRetryQuery, tweetID)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}
