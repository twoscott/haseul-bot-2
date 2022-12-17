package commanddb

import (
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
)

type Command struct {
	GuildID discord.GuildID `db:"guildid"`
	Name    string          `db:"name"`
	Content string          `db:"content"`
	Uses    int64           `db:"uses"`
	Created time.Time       `db:"created"`
}

const (
	createCommandsTableQuery = `
		CREATE TABLE IF NOT EXISTS Commands(
			guildID INT8         NOT NULL,
			name    VARCHAR(32)  NOT NULL,
			content VARCHAR(256) NOT NULL,
			uses    INT8         NOT NULL DEFAULT 0,
			created TIMESTAMPTZ  NOT NULL DEFAULT now(),
			PRIMARY KEY(guildID, name)
		)`
	addCommandQuery = `
		INSERT INTO Commands VALUES($1, $2, $3) ON CONFLICT DO NOTHING`
	getCommandQuery = `
		SELECT * FROM Commands WHERE guildID = $1 AND name = $2`
	getContentQuery = `
		SELECT content FROM Commands WHERE guildID = $1 AND name = $2`
	getAllCommandsByGuildQuery = `
		SELECT * FROM Commands WHERE guildID = $1`
	deleteCommandQuery = `
		DELETE FROM Commands WHERE guildID = $1 AND name = $2`
	useCommandQuery = `
		UPDATE Commands SET uses = uses + 1 WHERE guildID = $1 AND name = $2`
)

// Add adds a custom server command to a server.
func (db *DB) Add(guildID discord.GuildID, name, content string) (bool, error) {
	res, err := db.Exec(addCommandQuery, guildID, name, content)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// GetCommand returns the custom command with the provided name belonging to
// the provided guild.
func (db *DB) GetCommand(guildID discord.GuildID, name string) (*Command, error) {
	var command Command
	err := db.Get(&command, getCommandQuery, guildID, name)

	return &command, err
}

// GetContent gets the content for a named command.
func (db *DB) GetContent(guildID discord.GuildID, name string) (string, error) {
	var content string
	err := db.Get(&content, getContentQuery, guildID, name)

	return content, err
}

// GetAllByGuild returns all the custom commands for a given server.
func (db *DB) GetAllByGuild(guildID discord.GuildID) ([]Command, error) {
	var commands []Command
	err := db.Select(&commands, getAllCommandsByGuildQuery, guildID)

	return commands, err
}

// Delete deletes a custom server command from a server.
func (db *DB) Delete(guildID discord.GuildID, name string) (bool, error) {
	res, err := db.Exec(deleteCommandQuery, guildID, name)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// Use increments the uses value for a command.
func (db *DB) Use(guildID discord.GuildID, name string) (bool, error) {
	res, err := db.Exec(useCommandQuery, guildID, name)
	if err != nil {
		return false, err
	}

	used, err := res.RowsAffected()
	return used > 0, err
}
