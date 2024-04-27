package invitedb

import "github.com/diamondburned/arikawa/v3/discord"

// Invite represents a partial Discord invite used for guild invite tracking.
type Invite struct {
	Code    string
	GuildID discord.GuildID
	Uses    int
}

const (
	createGuildInvitesTableQuery = `
		CREATE TABLE IF NOT EXISTS GuildInvites(
			code    VARCHAR(32) NOT NULL,
			guildID INT8        NOT NULL,
			uses    INT         NOT NULL,
			PRIMARY KEY(code)
		)`

	addOrUpdateInviteQuery = `
		INSERT INTO GuildInvites VALUES($1, $2, $3)
		ON CONFLICT(code) DO UPDATE SET uses = $3`
	removeInviteQuery = `
		DELETE FROM GuildInvites WHERE code = $1`
	getGuildInvitesQuery = `
		SELECT * FROM GuildInvites WHERE guildID = $1`
)

// Add adds an invite to be tracked by the database, or if it exsits,
// updates the uses field to the new amount.
func (db *DB) Add(code string, guildID discord.GuildID, uses int) error {
	_, err := db.Exec(addOrUpdateInviteQuery, code, guildID, uses)
	return err
}

// AddAll adds multiple invites to be tracked by the database, or if it exsits,
// updates a uses field to the new amount.
func (db *DB) AddAll(guildID discord.GuildID, invites []discord.Invite) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, inv := range invites {
		tx.Exec(addOrUpdateInviteQuery, inv.Code, guildID, inv.Uses)
	}

	return tx.Commit()
}

// Remove removes an invite from being tracked in the database.
func (db *DB) Remove(code string) (bool, error) {

	res, err := db.Exec(removeInviteQuery, code)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// GetAllByGuild returns the tracked invites for a provided guild ID.
func (db *DB) GetAllByGuild(guildID discord.GuildID) ([]Invite, error) {
	var invites []Invite
	err := db.Select(&invites, getGuildInvitesQuery, guildID)

	return invites, err
}
