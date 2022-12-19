package reminderdb

import (
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
)

type Reminder struct {
	ID      int32          `db:"id"`
	UserID  discord.UserID `db:"userid"`
	Time    time.Time      `db:"time"`
	Content string         `db:"content"`
	Created time.Time      `db:"created"`
}

const (
	createRemindersTableQuery = `
		CREATE TABLE IF NOT EXISTS Reminders(
			ID      SERIAL,
			userID  INT8          NOT NULL,
			time    TIMESTAMP     NOT NULL,
			content VARCHAR(2048) NOT NULL,
			created TIMESTAMPTZ   NOT NULL DEFAULT now(),
			PRIMARY KEY(ID)
		)`
	addReminderQuery = `
		INSERT INTO Reminders (userID, time, content) 
		VALUES($1, $2, $3)
		RETURNING ID`
	deleteReminderQuery = `
		DELETE FROM Reminders WHERE userID = $1 AND ID = $2`
	clearRemindersQuery   = `DELETE FROM Reminders WHERE userID = $1`
	getAllRemindersByUser = `SELECT * FROM Reminders WHERE userID = $1`
	getOverdueReminders   = `SELECT * FROM Reminders WHERE time <= now()`
)

// Add adds a reminder for a user.
func (db *DB) Add(
	userID discord.UserID, time time.Time, content string) (int32, error) {

	var id int32
	err := db.Get(&id, addReminderQuery, userID, time, content)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// DeleteForUser deletes a reminder for a user.
func (db *DB) DeleteForUser(userID discord.UserID, id int32) (bool, error) {
	res, err := db.Exec(deleteReminderQuery, userID, id)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// ClearByUser deletes all reminders for a user.
func (db *DB) ClearByUser(userID discord.UserID) (int64, error) {
	res, err := db.Exec(clearRemindersQuery, userID)
	if err != nil {
		return 0, err
	}

	deleted, err := res.RowsAffected()
	return deleted, err
}

// GetAllByUser returns all the reminders for a user.
func (db *DB) GetAllByUser(userID discord.UserID) ([]Reminder, error) {
	var reminders []Reminder
	err := db.Select(&reminders, getAllRemindersByUser, userID)

	return reminders, err
}

// GetOverdueReminders returns all reminders that are ready to be sent to users.
func (db *DB) GetOverdueReminders() ([]Reminder, error) {
	var reminders []Reminder
	err := db.Select(&reminders, getOverdueReminders)

	return reminders, err
}
