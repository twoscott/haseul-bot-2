package repdb

import (
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
)

const maxRepsRemaining = 3

const (
	createRepHistoryTableQuery = `
		CREATE TABLE IF NOT EXISTS RepHistory(
			senderID   INT8        NOT NULL,
			receiverID INT8        NOT NULL,
			time       TIMESTAMPTZ NOT NULL DEFAULT now(),
			CHECK (senderID <> receiverID), 
			PRIMARY KEY (senderID, receiverID)
		)`
	addHistoryEntryQuery = `
		INSERT INTO RepHistory VALUES ($1, $2)
		ON CONFLICT (senderID, receiverID) DO
		UPDATE SET time = now()`
	getUserRepsRemainingQuery = `
		SELECT $1 - COUNT(*) FROM RepHistory 
		WHERE senderID = $2 AND time::DATE = CURRENT_DATE`
	getUserRepTimeQuery = `
		SELECT time FROM RepHistory
		WHERE senderID = $1 AND receiverID = $2`
)

// GetUserRepsRemaining returns the number of reps remaining for a user
func (db *DB) GetUserRepsRemaining(
	userID discord.UserID) (remaining int64, err error) {

	return remaining, db.Get(
		&remaining,
		getUserRepsRemainingQuery,
		maxRepsRemaining,
		userID,
	)
}

// GetUserLastRepTime returns the time when a user last gave a rep to someone.
func (db *DB) GetUserLastRepTime(
	senderID, receiverID discord.UserID) (t time.Time, err error) {

	return t, db.Get(
		&t,
		getUserRepTimeQuery,
		senderID,
		receiverID,
	)
}
