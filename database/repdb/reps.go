package repdb

import (
	"errors"

	"github.com/diamondburned/arikawa/v3/discord"
)

type RepUser struct {
	UserID discord.UserID `db:"userid"`
	Rep    int            `db:"rep"`
}

const (
	createRepTableQuery = `
		CREATE TABLE IF NOT EXISTS UserRep(
			userID    INT8 NOT NULL,
			rep       INT  NOT NULL,
			PRIMARY KEY(userID)
		)`
	repUserQuery = `
		INSERT INTO UserRep VALUES($1, 1)
		ON CONFLICT(userID) DO UPDATE SET rep = UserRep.rep + 1
		RETURNING rep`
	getUserXPQuery = `
		SELECT rep FROM UserRep WHERE userID = $1`
	getTopUsersQuery = `
		SELECT * FROM UserRep
		ORDER BY rep DESC
		LIMIT $1`
	getAllUsersQuery = `
		SELECT * FROM UserRep`
	getEntriesSizeQuery = `
		SELECT SUM(rep) FROM UserRep`
)

// RepUser adds a rep to a user.
func (db *DB) RepUser(senderID, targetID discord.UserID) (rep int, err error) {
	if senderID == targetID {
		return 0, errors.New("sender and target rep users cannot be the same")
	}

	tx, err := db.Beginx()
	if err != nil {
		return 0, err
	}

	defer tx.Rollback()

	err = tx.Get(&rep, repUserQuery, targetID)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(addHistoryEntryQuery, senderID, targetID)
	if err != nil {
		return 0, err
	}

	_, err = db.AddOrUpdateRepStreak(senderID, targetID)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return rep, nil
}

// GetUserRep returns the rep for a user.
func (db *DB) GetUserRep(userID discord.UserID) (rep int, err error) {
	return rep, db.Get(&rep, getUserXPQuery, userID)
}

// GetTopUsers returns the most repped users.
func (db *DB) GetTopUsers(limit int64) (users []RepUser, err error) {
	return users, db.Select(&users, getTopUsersQuery, limit)
}

// GetAllUsers returns all users and their rep scores.
func (db *DB) GetAllUsers() (users []RepUser, err error) {
	return users, db.Select(&users, getAllUsersQuery)
}

// GetTotalReps returns the total number of reps between all users.
func (db *DB) GetTotalReps() (reps int64, err error) {
	return reps, db.Get(&reps, getEntriesSizeQuery)
}
