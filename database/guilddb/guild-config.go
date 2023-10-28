package guilddb

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

// Config represents a guild config database entry.
type Config struct {
	GuildID              discord.GuildID   `db:"guildid"`
	LegacyPrefix         rune              `db:"legacyprefix"`
	MemberLogsChannelID  discord.ChannelID `db:"memberlogschannelid"`
	MessageLogsChannelID discord.ChannelID `db:"messagelogschannelid"`

	Welcome Welcome `db:""`
}

const (
	createGuildConfigsTableQuery = `
		CREATE TABLE IF NOT EXISTS GuildConfigs(
			guildID              INT8          NOT NULL,

			legacyPrefix         CHAR(1)       DEFAULT '.',
			memberLogsChannelID  INT8          DEFAULT 0,
			messageLogsChannelID INT8          DEFAULT 0,

			welcomeChannelID     INT8          DEFAULT 0,
			welcomeTitle         VARCHAR(32)   DEFAULT '',
			welcomeMessage       VARCHAR(1024) DEFAULT '',
			welcomeColour        INT4		   DEFAULT NULL,
			
			PRIMARY KEY(guildID)
		)`

	getConfigsQuery = `SELECT * FROM GuildConfigs`
	getConfigQuery  = `SELECT * FROM GuildConfigs WHERE guildID = $1`
	addConfigQuery  = `
		INSERT INTO GuildConfigs(guildID) VALUES($1) ON CONFLICT DO NOTHING`
	getLegacyPrefixQuery = `
		SELECT legacyPrefix FROM GuildConfigs WHERE guildID = $1`
	setMemberLogsChannelQuery = `
		UPDATE GuildConfigs SET memberLogsChannelID = $1
		WHERE guildID = $2`
	setMemberLogsChannelNullQuery = `
		UPDATE GuildConfigs SET memberLogsChannelID = 0
		WHERE guildID = $1`
	getMemberLogsQuery = `
		SELECT memberLogsChannelID FROM GuildConfigs WHERE guildID = $1`
	setMessageLogsChannelQuery = `
		UPDATE GuildConfigs SET messageLogsChannelID = $1
		WHERE guildID = $2`
	setMessageLogsChannelNullQuery = `
		UPDATE GuildConfigs SET messageLogsChannelID = 0
		WHERE guildID = $1`
	getMessageLogsQuery = `
		SELECT messageLogsChannelID FROM GuildConfigs WHERE guildID = $1`
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

// GetLegacyPrefix returns the legacy prefix for the guild - the prefix uses
// for old commands.
func (db *DB) GetLegacyPrefix(
	guildID discord.GuildID) (prefix string, err error) {

	return prefix, db.Get(&prefix, getLegacyPrefixQuery, guildID)
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

// GetMemberLogsChannel returns the member logs channelID of the guild config
// corresponding to the provided guild ID.
func (db *DB) GetMemberLogsChannel(
	guildID discord.GuildID) (discord.ChannelID, error) {

	var id discord.ChannelID
	err := db.Get(&id, getMemberLogsQuery, guildID)

	return id, err
}

// SetMessageLogsChannel updates the guild config of the given guild ID and sets
// the message logs channel ID to the provided channel ID.
func (db *DB) SetMessageLogsChannel(
	guildID discord.GuildID, channelID discord.ChannelID) (bool, error) {

	res, err := db.Exec(setMessageLogsChannelQuery, channelID, guildID)
	if err != nil {
		return false, err
	}

	updated, err := res.RowsAffected()
	return updated > 0, err
}

// DisableMessageLogs updates the guild config of the given guild ID and sets
// the message logs channel ID to a value that can be interpreted as null.
func (db *DB) DisableMessageLogs(guildID discord.GuildID) (bool, error) {
	res, err := db.Exec(setMessageLogsChannelNullQuery, guildID)
	if err != nil {
		return false, err
	}

	updated, err := res.RowsAffected()
	return updated > 0, err
}

// GetMessageLogsChannel returns the message logs channelID of the guild config
// corresponding to the provided guild ID.
func (db *DB) GetMessageLogsChannel(
	guildID discord.GuildID) (discord.ChannelID, error) {

	var id discord.ChannelID
	err := db.Get(&id, getMessageLogsQuery, guildID)

	return id, err
}
