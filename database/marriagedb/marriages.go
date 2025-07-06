package marriagedb

import (
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
)

type Marriage struct {
	SpouseID1 discord.UserID `db:"spouseid1"`
	SpouseID2 discord.UserID `db:"spouseid2"`
	MarriedAt time.Time      `db:"marriedat"`
}

// Spouse returns the spouse of a user.
func (m Marriage) Spouse(user discord.UserID) discord.UserID {
	if user == m.SpouseID1 {
		return m.SpouseID2
	} else if user == m.SpouseID2 {
		return m.SpouseID1
	}
	return discord.NullUserID
}

const (
	createMarriagesTableQuery = `
		CREATE TABLE IF NOT EXISTS Marriages(
			spouseID1 INT8 		  NOT NULL,
			spouseID2 INT8 		  NOT NULL,
			marriedAt TIMESTAMPTZ NOT NULL DEFAULT now(),
			PRIMARY KEY(spouseID1, spouseID2)
		)`
	addMarriageQuery = `
		INSERT INTO Marriages (spouseID1, spouseID2) 
		VALUES($1, $2)
		ON CONFLICT DO NOTHING`
	removeMarriageQuery = `
		DELETE FROM Marriages 
		WHERE spouseID1 = $1 OR spouseID2 = $1`
	getUserMarriageQuery = `
		SELECT * FROM Marriages 
		WHERE spouseID1 = $1 OR spouseID2 = $1`
)

// Add adds a marriage between two users.
func (db *DB) Add(spouseID1, spouseID2 discord.UserID) (bool, error) {
	res, err := db.Exec(addMarriageQuery, spouseID1, spouseID2)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// Remove removes a marriage between two users.
func (db *DB) Remove(spouse discord.UserID) (int64, error) {
	res, err := db.Exec(removeMarriageQuery, spouse)
	if err != nil {
		return 0, err
	}

	deleted, err := res.RowsAffected()
	return deleted, err
}

// GetUserMarriage returns the marriage of a user.
func (db *DB) GetUserMarriage(userID discord.UserID) (Marriage, error) {
	var marriage Marriage
	err := db.Get(&marriage, getUserMarriageQuery, userID)
	return marriage, err
}
