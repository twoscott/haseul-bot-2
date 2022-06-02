package twitterdb

import "github.com/diamondburned/arikawa/v3/discord"

// Mention represents a Twitter mention database entry.
type Mention struct {
	ChannelID     discord.ChannelID `db:"channelid"`
	TwitterUserID int64             `db:"twitteruserid"`
	RoleID        discord.RoleID    `db:"roleid"`
}

const (
	createTwitterMentionsTableQuery = `
		CREATE TABLE IF NOT EXISTS TwitterMentions(
			channelID     INT8 NOT NULL,
			twitterUserID INT8 NOT NULL,
			roleID        INT8 NOT NULL,
			PRIMARY KEY(channelID, twitterUserID, roleID),
			FOREIGN KEY(channelID, twitterUserID) 
			REFERENCES TwitterFeeds(channelID, twitterUserID)
		)`

	addMentionQuery = `
		INSERT INTO TwitterMentions VALUES($1, $2, $3) ON CONFLICT DO NOTHING`
	getMentionRolesQuery = `
		SELECT roleID FROM TwitterMentions 
		WHERE channelID = $1 AND twitterUserID = $2`
	getMentionsQuery = `
			SELECT * FROM TwitterMentions 
			WHERE channelID = $1 AND twitterUserID = $2`
	getMentionsByGuildQuery = `
		SELECT * FROM TwitterMentions WHERE channelID IN (
			SELECT channelID FROM TwitterFeeds WHERE guildID = $1
		)`
	removeMentionQuery = `
		DELETE FROM TwitterMentions 
		WHERE channelID = $1 AND twitterUserID = $2 AND roleID = $3`
	removeMentionsQuery = `
		DELETE FROM TwitterMentions WHERE channelID = $1 AND twitterUserID = $2`
	clearGuildMentionsQuery = `
		DELETE FROM TwitterMentions WHERE channelID IN (
			SELECT channelID FROM TwitterFeeds WHERE guildID = $1
		)`
)

// AddMention adds a Twitter mention to the database.
func (db *DB) AddMention(channelID discord.ChannelID, twitterUserID int64, roleID discord.RoleID) (bool, error) {
	res, err := db.Exec(addMentionQuery, channelID, twitterUserID, roleID)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// GetMentionRoles returns all Twitter mentions for a Twitter feed.
func (db *DB) GetMentionRoles(channelID discord.ChannelID, twitterUserID int64) ([]discord.RoleID, error) {
	var roleIDs []discord.RoleID
	err := db.Select(&roleIDs, getMentionRolesQuery, channelID, twitterUserID)

	return roleIDs, err
}

// GetMentionsByGuild returns all Twitter mentions in a guild ID.
func (db *DB) GetMentions(channelID discord.ChannelID, twitterUserID int64) ([]Mention, error) {
	var mentionRoles []Mention
	err := db.Select(&mentionRoles, getMentionsQuery, channelID, twitterUserID)

	return mentionRoles, err
}

// GetMentionsByGuild returns all Twitter mentions in a guild ID.
func (db *DB) GetMentionsByGuild(guildID discord.GuildID) ([]Mention, error) {
	var mentionRoles []Mention
	err := db.Select(&mentionRoles, getMentionsByGuildQuery, guildID)

	return mentionRoles, err
}

// RemoveMention removes a Twitter mention.
func (db *DB) RemoveMention(
	channelID discord.ChannelID, twitterUserID int64, roleID discord.RoleID) (bool, error) {
	res, err := db.Exec(removeMentionQuery, channelID, twitterUserID, roleID)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// RemoveMentions removes all Twitter mentions for a Twitter feed.
func (db *DB) RemoveMentions(
	channelID discord.ChannelID, twitterUserID int64) (bool, error) {
	res, err := db.Exec(removeMentionsQuery, channelID, twitterUserID)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// ClearGuildMentions removes all Twitter mentions in a guild ID.
func (db *DB) ClearGuildMentions(guildID discord.GuildID) (int64, error) {
	res, err := db.Exec(clearGuildMentionsQuery, guildID)
	if err != nil {
		return 0, err
	}

	cleared, err := res.RowsAffected()
	return cleared, err
}
