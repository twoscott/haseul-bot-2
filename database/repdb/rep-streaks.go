package repdb

import (
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
)

type RepStreak struct {
	UserID1  discord.UserID `db:"userid1"`
	UserID2  discord.UserID `db:"userid2"`
	FirstRep time.Time      `db:"firstrep"`
}

const (
	createRepStreaksTableQuery = `
		CREATE TABLE IF NOT EXISTS RepStreaks(
			userID1  INT8        NOT NULL,
			userID2  INT8        NOT NULL,
			firstRep TIMESTAMPTZ NOT NULL DEFAULT now(),
			CHECK (userID1 <> userID2),
			CHECK (userID1 < userID2),
			PRIMARY KEY (userID1, userID2)
		)`
	addOrUpdateRepStreakQuery = `
		INSERT INTO RepStreaks VALUES($1, $2)
		ON CONFLICT(userID1, userID2) DO
		UPDATE SET firstRep = CASE
			WHEN (
				SELECT COUNT(*) FROM RepHistory AS rh
				WHERE rh.senderID IN ($1, $2) 
					AND rh.receiverID IN ($1, $2)
					AND now() - rh.time < INTERVAL '36 hours'
			) = 0 THEN now()
			ELSE RepStreaks.firstRep
		END`
	updateRepStreaksQuery = `
		DELETE FROM RepStreaks AS rs
			WHERE (
				SELECT COUNT(*) FROM RepHistory AS rh
				WHERE rh.senderID IN (rs.userID1, rs.userID2) 
					AND rh.receiverID IN (rs.userID1, rs.userID2)
					AND now() - rh.time < INTERVAL '36 hours'
			) < 2`
	getUserStreakQuery = `
		SELECT CURRENT_DATE - firstRep::DATE FROM RepStreaks
		WHERE userID1 IN ($1, $2) AND userID2 IN ($1, $2)`
	getTopStreaksQuery = `
		SELECT * FROM RepStreaks WHERE now() - firstRep > INTERVAL '24 hours'
		ORDER BY firstRep ASC
		LIMIT $1`
	getEntriesSizeQuery = `
		SELECT COUNT(*) FROM RepStreaks 
		WHERE now() - firstRep > INTERVAL '24 hours'`
)

// AddOrUpdateRepStreak adds a rep streak start entry to the database if it
// doesn't exist, or updates the time according to whether the rep streak should
// be reset or continue, based on the users' rep history.
func (db *DB) AddOrUpdateRepStreak(
	senderID, targetID discord.UserID) (bool, error) {

	var userID1, userID2 discord.UserID
	if senderID < targetID {
		userID1, userID2 = senderID, targetID
	} else {
		userID1, userID2 = targetID, senderID
	}

	res, err := db.Exec(addOrUpdateRepStreakQuery, userID1, userID2)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// UpdateRepStreaks clears any rep streaks that have fallen past the max rep
// age time.
func (db *DB) UpdateRepStreaks() (int64, error) {
	res, err := db.Exec(updateRepStreaksQuery)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// GetUserStreak returns the number of days a rep streak is at between two
// users.
func (db *DB) GetUserStreak(
	userID1, userID2 discord.UserID) (streak int, err error) {

	return streak, db.Get(&streak, getUserStreakQuery, userID1, userID2)
}

// GetTopStreaks returns the pairs of users with the longest rep streaks.
func (db *DB) GetTopStreaks(limit int64) (streaks []RepStreak, err error) {
	return streaks, db.Select(&streaks, getTopStreaksQuery, limit)
}

// GetTotalStreaks returns the number of ongoing streaks.
func (db *DB) GetTotalStreaks() (count int, err error) {
	return count, db.Get(&count, getEntriesSizeQuery)
}
