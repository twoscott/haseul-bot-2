package youtubedb

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
)

// HistoryEntry represents an entry in a user's YouTube search history.
type HistoryEntry struct {
	UserID        discord.UserID        `db:"userid"`
	InteractionID discord.InteractionID `db:"interactionid"`
	Query         string                `db:"query"`
}

const historyEntriesToKeep = 25

const (
	createYouTubeHistoryTable = `
		CREATE TABLE IF NOT EXISTS YouTubeHistory(
			userID        INT8 NOT NULL,
			interactionID INT8 NOT NULL,
			query		  TEXT NOT NULL,
			PRIMARY KEY(userID, interactionID)
		)`

	addHistoryEntryQuery = `
		INSERT INTO YouTubeHistory VALUES($1, $2, $3)`

	deleteOldHistoryQuery = `
		DELETE FROM YouTubeHistory 
		WHERE interactionID IN (
			SELECT interactionID FROM YouTubeHistory
			WHERE userID = $1
			ORDER BY interactionID DESC
			OFFSET $2
		)`

	clearHistoryQuery = `
		DELETE FROM YouTubeHistory WHERE userID = $1`

	getHistoryQuery = `
		SELECT query FROM YouTubeHistory 
		WHERE userID = $1
		ORDER BY interactionID DESC`
)

// AddHistoryAndClear adds a new YouTube search to a user's search history and
// clears old searches, keeping only the most recent 10 searches in history.
func (db *DB) AddHistoryAndClear(
	userID discord.UserID,
	interactionID discord.InteractionID,
	query string) error {

	isEnabled, err := db.GetHistoryToggle(userID)
	if err != nil {
		return err
	}
	if !isEnabled {
		return nil
	}

	err = db.addHistoryEntry(userID, interactionID, query)
	if err != nil {
		return err
	}

	_, err = db.deleteOldHistory(userID)
	return err
}

// addHistoryEntry adds a new YouTube search to a user's search history and
// clears old searches, keeping only the most recent 10 searches in history.
func (db *DB) addHistoryEntry(
	userID discord.UserID,
	interactionID discord.InteractionID,
	query string) error {

	_, err := db.Exec(addHistoryEntryQuery, userID, interactionID, query)
	return err
}

// deleteOldHistory deletes old history entries past the most recent
// 10 searches.
func (db *DB) deleteOldHistory(userID discord.UserID) (int64, error) {
	res, err := db.Exec(deleteOldHistoryQuery, userID, historyEntriesToKeep)
	if err != nil {
		return 0, err
	}

	deleted, err := res.RowsAffected()
	log.Printf("%d deleted from history", deleted)
	return deleted, err
}

// ClearHistory deletes all YouTube search history for a user.
func (db *DB) ClearHistory(userID discord.UserID) (int64, error) {
	res, err := db.Exec(clearHistoryQuery, userID)
	if err != nil {
		return 0, err
	}

	cleared, err := res.RowsAffected()
	return cleared, err
}

// GetHistory returns the stored YouTube search history for a user.
func (db *DB) GetHistory(userID discord.UserID) ([]string, error) {
	var entries []HistoryEntry
	err := db.Select(&entries, getHistoryQuery, userID)

	queries := make([]string, 0, len(entries))
	for _, e := range entries {
		queries = append(queries, e.Query)
	}

	return queries, err
}
