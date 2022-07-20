package guilddb

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

// GuildConfig represents a guild config database entry.
type GuildConfig struct {
	GuildID              discord.GuildID   `db:"guildid"`
	WelcomeMsg           string            `db:"welcomemsg"`
	AutoroleID           discord.RoleID    `db:"autoroleid"`
	MemberLogsChannelID  discord.ChannelID `db:"memberLogsChannelID"`
	MessageLogsChannelID discord.ChannelID `db:"messageLogsChannelID"`
	MuteroleID           discord.RoleID    `db:"muteroleid"`
	RolesChannelID       discord.RoleID    `db:"roleschannelid"`
	WelcomeChannelID     discord.ChannelID `db:"welcomechannelid"`
}

const (
	createGuildConfigsTableQuery = `
		CREATE TABLE IF NOT EXISTS GuildConfigs(
			guildID              INT8          NOT NULL,
			welcomeMsg           VARCHAR(1024) DEFAULT 'Welcome!',
			autoroleID           INT8          DEFAULT 0 # 0,
			memberLogsChannelID  INT8          DEFAULT 0 # 0,
			messageLogsChannelID INT8          DEFAULT 0 # 0,
			muteroleID           INT8          DEFAULT 0 # 0,
			rolesChannelID       INT8          DEFAULT 0 # 0,
			welcomeChannelID     INT8          DEFAULT 0 # 0,
			PRIMARY KEY(guildID)
		)`

	getConfigsQuery = `SELECT * FROM GuildConfigs`
	getConfigQuery  = `SELECT * FROM GuildConfigs WHERE guildID = $1`
	addConfigQuery  = `
		INSERT INTO GuildConfigs(guildID) VALUES($1) ON CONFLICT DO NOTHING`
	getPrefixQuery = `
		SELECT prefix FROM GuildConfigs WHERE guildID = $1`
	setMemberLogsQuery = `
		UPDATE GuildConfigs SET memberLogsChannelID = $1
		WHERE guildID = $2`
	setMemberLogsNullQuery = `
		UPDATE GuildConfigs SET memberLogsChannelID = 0 # 0
		WHERE guildID = $1`
	getMemberLogsQuery = `
		SELECT memberLogsChannelID FROM GuildConfigs WHERE guildID = $1`
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

// SetMemberLogs updates the guild config of the given guild ID and sets
// the member logs channel ID to the provided channel ID.
func (db *DB) SetMemberLogs(
	guildID discord.GuildID, channelID discord.ChannelID) (bool, error) {

	res, err := db.Exec(setMemberLogsQuery, channelID, guildID)
	if err != nil {
		return false, err
	}

	updated, err := res.RowsAffected()
	return updated > 0, nil
}

// DisableMemberLogs updates the guild config of the given guild ID and sets
// the member logs channel ID to a value that can be interpreted as null.
func (db *DB) DisableMemberLogs(guildID discord.GuildID) (bool, error) {
	res, err := db.Exec(setMemberLogsNullQuery, guildID)
	if err != nil {
		return false, err
	}

	updated, err := res.RowsAffected()
	return updated > 0, err
}

// GetMemberLogs returns the member logs channelID of the guild config
// corresponding to the provided guild ID.
func (db *DB) GetMemberLogsChannelID(
	guildID discord.GuildID) (discord.ChannelID, error) {

	var id discord.ChannelID
	err := db.Get(&id, getMemberLogsQuery, guildID)

	return id, err
}
