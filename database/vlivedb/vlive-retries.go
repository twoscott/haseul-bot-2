package vlivedb

const (
	createVLIVERetriesTableQuery = `
		CREATE TABLE IF NOT EXISTS VLIVERetries(
			postID VARCHAR(64) NOT NULL,
			PRIMARY KEY(postID)
		)`

	addRetryQuery = `
		INSERT INTO VLIVERetries VALUES($1) ON CONFLICT DO NOTHING`
	getAllRetriesQuery = `
		SELECT postID FROM VLIVERetries`
	removeRetryQuery = `
		DELETE FROM VLIVERetries WHERE postID = $1`
)

// AddRetry returns adds a VLIVE post to retry sending to the database.
func (db *DB) AddRetry(postID string) (bool, error) {
	res, err := db.Exec(addRetryQuery, postID)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// GetAllUsers returns all VLIVE posts that previously failed to send from
// the database.
func (db *DB) GetAllRetries() ([]string, error) {
	var postIDs []string
	err := db.Select(&postIDs, getAllRetriesQuery)

	return postIDs, err
}

// RemoveRetry removes a retry VLIVE post from the database.
func (db *DB) RemoveRetry(postID string) (bool, error) {
	res, err := db.Exec(removeRetryQuery, postID)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}
