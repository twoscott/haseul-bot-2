package guilddb

import (
	"database/sql"

	"github.com/diamondburned/arikawa/v3/discord"
)

// GuildConfig represents a guild config database entry.
type GuildConfig struct {
	GuildID           discord.GuildID `db:"guildid"`
	Prefix            string          `db:"prefix"`
	WelcomeMsg        string          `db:"welcomemsg"`
	AutoroleOn        bool            `db:"autoroleon"`
	CommandsOn        bool            `db:"commandson"`
	JoinLogsOn        bool            `db:"joinlogson"`
	MsgLogsOn         bool            `db:"msglogson"`
	PollOn            bool            `db:"pollon"`
	RolesOn           bool            `db:"roleson"`
	WelcomeOn         bool            `db:"welcomeon"`
	AutoroleID        sql.NullInt64   `db:"autoroleid"`
	JoinLogsChannelID sql.NullInt64   `db:"joinlogschannelid"`
	MsgLogsChannelID  sql.NullInt64   `db:"msglogschannelid"`
	MuteroleID        sql.NullInt64   `db:"muteroleid"`
	RolesChannelID    sql.NullInt64   `db:"roleschannelid"`
	WelcomeChannelID  sql.NullInt64   `db:"welcomechannelid"`
}

const (
	createGuildConfigsTableQuery = `
		CREATE TABLE IF NOT EXISTS GuildConfig(
			guildID           INT8       NOT NULL,
			prefix            VARCHAR(3) DEFAULT '.',
			welcomeMsg        TEXT       DEFAULT 'Welcome!',
			autoroleOn        BOOLEAN    DEFAULT FALSE,
			commandsOn        BOOLEAN    DEFAULT TRUE, 
			joinLogsOn        BOOLEAN    DEFAULT FALSE, 
			msgLogsOn         BOOLEAN    DEFAULT FALSE,
			pollOn            BOOLEAN    DEFAULT FALSE,
			rolesOn           BOOLEAN    DEFAULT FALSE,
			welcomeOn         BOOLEAN    DEFAULT FALSE,
			autoroleID        INT8  ,
			joinLogsChannelID INT8  ,
			msgLogsChannelID  INT8  ,
			muteroleID        INT8  ,
			rolesChannelID    INT8  ,
			welcomeChannelID  INT8  ,
			PRIMARY KEY(guildID)
		)`

	getConfigsQuery = `SELECT * FROM GuildConfig`
	getConfigQuery  = `SELECT * FROM GuildConfig WHERE guildID = $1`
	addConfigQuery  = `
		INSERT INTO GuildConfig(guildID) VALUES($1) ON CONFLICT DO NOTHING`
	getPrefixQuery = `
		SELECT prefix FROM GuildConfig WHERE guildID = $1`
)

// GetConfigs returns all guild configs from the database.
func (db *DB) GetConfigs() ([]GuildConfig, error) {
	var configs []GuildConfig
	err := db.Select(&configs, getConfigsQuery)
	return configs, err
}

// GetConfig returns a guild config for the given guild ID.
func (db *DB) GetConfig(guildID discord.GuildID) (*GuildConfig, error) {
	var config GuildConfig
	err := db.Get(&config, getConfigQuery, guildID)
	return &config, err
}

// Add adds a guild config for the given guild ID to the database.
func (db *DB) Add(guildID discord.GuildID) (bool, error) {
	res, err := db.Exec(addConfigQuery, guildID)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 1, err
}
