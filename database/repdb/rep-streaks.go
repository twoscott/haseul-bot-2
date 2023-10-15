package repdb

import "github.com/diamondburned/arikawa/v3/discord"

// type RepUser struct {
// 	UserID discord.UserID `db:"userid"`
// 	Rep    int            `db:"rep"`
// }

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
				SELECT COUNT(*) FROM RepHistory
				WHERE RepHistory.senderID IN ($1, $2) 
					AND RepHistory.receiverID IN ($1, $2)
					AND now() - RepHistory.time < INTERVAL '36 hours'
			) = 0 THEN now()
			ELSE RepStreaks.firstRep
		END`
	updateRepStreaksQuery = `
		DELETE rs FROM RepStreaks AS rs
		INNER JOIN RepHistory AS rh
			ON now() - rh.time <= INTERVAL '36 hours'
			AND rh.senderID IN (rs.userID1, rs.userID2)
			AND rh.receiverID IN (rs.userID1, rs.userID2)`
	getUserStreakQuery = `
		SELECT CURRENT_DATE - firstRep::DATE FROM RepStreaks
		WHERE userID1 IN ($1, $2) AND userID2 IN ($1, $2)`
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
