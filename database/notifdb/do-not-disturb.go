package notifdb

import "github.com/diamondburned/arikawa/v3/discord"

const (
	createNotiDnDTableQuery = `
		CREATE TABLE IF NOT EXISTS NotiDnD(
			userID    INT8 NOT NULL,
			PRIMARY KEY(userID)
		)`
	addDnD = `
		INSERT INTO NotiDnD VALUES($1) ON CONFLICT DO NOTHING`
	removeDnD = `
		DELETE FROM NotiDnD WHERE userID = $1`
)

func (db *DB) ToggleDnD(userID discord.UserID) (bool, error) {
	res, err := db.Exec(addDnD, userID)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	if err != nil {
		return added > 0, err
	}
	if added < 1 {
		_, err = db.Exec(removeDnD, userID)
	}

	return added > 0, err
}
