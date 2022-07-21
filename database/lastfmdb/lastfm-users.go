package lastfmdb

import "github.com/diamondburned/arikawa/v3/discord"

const (
	createLastFmUsersTableQuery = `
		CREATE TABLE IF NOT EXISTS LastFmUsers(
			userID INT8        NOT NULL PRIMARY KEY,
			lfUser VARCHAR(15) NOT NULL
		)`

	setUpdateUserQuery = `
		INSERT INTO LastFmUsers VALUES($1, $2) 
		ON CONFLICT(userID) DO UPDATE SET lfUser = $2`
	deleteUserQuery = `DELETE FROM LastFmUsers WHERE userID = $1`
	getUserQuery    = `SELECT lfUser FROM LastFmUsers WHERE userID = $1`
)

// SetUser either adds a new entry to the database with the given
// user ID and Last.fm username, or updates the user ID's Last.fm username.
func (db *DB) SetUser(userID discord.UserID, lfUser string) error {
	_, err := db.Exec(setUpdateUserQuery, userID, lfUser)
	return err
}

// DeleteUser removes the Last.fm username for a given user ID.
func (db *DB) DeleteUser(userID discord.UserID) (bool, error) {
	res, err := db.Exec(deleteUserQuery, userID)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return deleted > 0, err
}

// GetUser gets the Last.fm username for a given user ID.
func (db *DB) GetUser(userID discord.UserID) (string, error) {
	var lfUser string
	err := db.Get(&lfUser, getUserQuery, userID)

	return lfUser, err
}
