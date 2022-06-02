package youtubedb

import (
	"database/sql"
	"errors"

	"github.com/diamondburned/arikawa/v3/discord"
)

// HistoryDisabled represents a user's history disabled option.
type HistoryDisabled struct {
	UserID discord.UserID `db:"userid"`
}

const (
	createHistoryToggleTableQuery = `
		CREATE TABLE IF NOT EXISTS YouTubeHistoryDisabled(
			userID INT8    NOT NULL,
			PRIMARY KEY(userID)
		)`

	disableHistory = `
		INSERT INTO YouTubeHistoryDisabled 
		VALUES($1) ON CONFLICT DO NOTHING`

	enableHistory = `
		DELETE FROM YouTubeHistoryDisabled WHERE userID = $1`

	getHistoryToggleQuery = `
		SELECT * FROM YouTubeHistoryDisabled WHERE userID = $1`
)

// ToggleHistory toggles whether a user's YouTube search history will be
// tracked or not.
func (db *DB) ToggleHistory(userID discord.UserID) error {
	res, err := db.Exec(disableHistory, userID)
	if err != nil {
		return err
	}

	disabled, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if disabled < 1 {
		_, err = db.Exec(enableHistory, userID)
	}

	return err
}

// GetHistoryToggle returns whether a user has YouTube search history enabled
// or disabled.
func (db *DB) GetHistoryToggle(userID discord.UserID) (bool, error) {
	var toggle HistoryDisabled
	err := db.Get(&toggle, getHistoryToggleQuery, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return true, nil
	}

	return !toggle.UserID.IsValid(), err
}
