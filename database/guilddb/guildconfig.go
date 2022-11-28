package guilddb

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

// Config represents a guild config database entry.
type Config struct {
	GuildID              discord.GuildID   `db:"guildid"`
	AutoroleID           discord.RoleID    `db:"autoroleid"`
	MemberLogsChannelID  discord.ChannelID `db:"memberlogschannelid"`
	MessageLogsChannelID discord.ChannelID `db:"messagelogschannelid"`
	MuteroleID           discord.RoleID    `db:"muteroleid"`
	RolesChannelID       discord.RoleID    `db:"roleschannelid"`

	Welcome Welcome `db:""`
}

const (
	createGuildConfigsTableQuery = `
		CREATE TABLE IF NOT EXISTS GuildConfigs(
			guildID              INT8          NOT NULL,
			autoroleID           INT8          DEFAULT 0 # 0,
			memberLogsChannelID  INT8          DEFAULT 0 # 0,
			messageLogsChannelID INT8          DEFAULT 0 # 0,
			muteroleID           INT8          DEFAULT 0 # 0,
			rolesChannelID       INT8          DEFAULT 0 # 0,
			welcomeChannelID     INT8          DEFAULT 0 # 0,
			welcomeTitle         VARCHAR(32)   DEFAULT '',
			welcomeMessage       VARCHAR(1024) DEFAULT '',
			welcomeColour        INT4		   DEFAULT NULL,
			PRIMARY KEY(guildID)
		)`

	getConfigsQuery = `SELECT * FROM GuildConfigs`
	getConfigQuery  = `SELECT * FROM GuildConfigs WHERE guildID = $1`
	addConfigQuery  = `
		INSERT INTO GuildConfigs(guildID) VALUES($1) ON CONFLICT DO NOTHING`
	getPrefixQuery = `
		SELECT prefix FROM GuildConfigs WHERE guildID = $1`
	setMemberLogsChannelQuery = `
		UPDATE GuildConfigs SET memberLogsChannelID = $1
		WHERE guildID = $2`
	setMemberLogsChannelNullQuery = `
		UPDATE GuildConfigs SET memberLogsChannelID = 0 # 0
		WHERE guildID = $1`
	getMemberLogsQuery = `
		SELECT memberLogsChannelID FROM GuildConfigs WHERE guildID = $1`
)

// Configs returns all guild configs from the database.
func (db *DB) Configs() ([]Config, error) {
	var configs []Config
	err := db.Select(&configs, getConfigsQuery)
	return configs, err
}

// Config returns a guild config for the given guild ID.
func (db *DB) Config(guildID discord.GuildID) (*Config, error) {
	var config Config
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

// SetMemberLogsChannel updates the guild config of the given guild ID and sets
// the member logs channel ID to the provided channel ID.
func (db *DB) SetMemberLogsChannel(
	guildID discord.GuildID, channelID discord.ChannelID) (bool, error) {

	res, err := db.Exec(setMemberLogsChannelQuery, channelID, guildID)
	if err != nil {
		return false, err
	}

	updated, err := res.RowsAffected()
	return updated > 0, err
}

// DisableMemberLogs updates the guild config of the given guild ID and sets
// the member logs channel ID to a value that can be interpreted as null.
func (db *DB) DisableMemberLogs(guildID discord.GuildID) (bool, error) {
	res, err := db.Exec(setMemberLogsChannelNullQuery, guildID)
	if err != nil {
		return false, err
	}

	updated, err := res.RowsAffected()
	return updated > 0, err
}

// MemberLogsChannel returns the greeting logs channelID of the guild config
// corresponding to the provided guild ID.
func (db *DB) MemberLogsChannel(
	guildID discord.GuildID) (discord.ChannelID, error) {

	var id discord.ChannelID
	err := db.Get(&id, getMemberLogsQuery, guildID)

	return id, err
}
