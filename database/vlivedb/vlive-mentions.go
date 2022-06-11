package vlivedb

import "github.com/diamondburned/arikawa/v3/discord"

// Mention represents a VLIVE mention database entry.
type Mention struct {
	BoardID   int64             `db:"boardid"`
	ChannelID discord.ChannelID `db:"channelid"`
	RoleID    discord.RoleID    `db:"roleid"`
}

const (
	createVLIVEMentionsTableQuery = `
		CREATE TABLE IF NOT EXISTS VLIVEMentions(
			boardID   INT8 NOT NULL,
			channelID INT8 NOT NULL,
			roleID    INT8 NOT NULL,
			PRIMARY KEY(boardID, channelID, roleID),
			FOREIGN KEY(boardID, channelID) 
			REFERENCES VLIVEFeeds(boardID, channelID)
		)`

	addMentionQuery = `
		INSERT INTO VLIVEMentions VALUES($1, $2, $3) ON CONFLICT DO NOTHING`
	getMentionRolesQuery = `
		SELECT roleID FROM VLIVEMentions 
		WHERE channelID = $1 AND boardID = $2`
	getMentionsQuery = `
			SELECT * FROM VLIVEMentions 
			WHERE channelID = $1 AND boardID = $2`
	getMentionsByGuildQuery = `
		SELECT * FROM VLIVEMentions WHERE channelID IN (
			SELECT channelID FROM VLIVEFeeds WHERE guildID = $1
		)`
	removeMentionQuery = `
		DELETE FROM VLIVEMentions 
		WHERE channelID = $1 AND boardID = $2 AND roleID = $3`
	removeMentionsQuery = `
		DELETE FROM VLIVEMentions WHERE channelID = $1 AND boardID = $2`
	clearGuildMentionsQuery = `
		DELETE FROM VLIVEMentions WHERE channelID IN (
			SELECT channelID FROM VLIVEFeeds WHERE guildID = $1
		)`
)

// AddMention adds a VLIVE mention to the database.
func (db *DB) AddMention(channelID discord.ChannelID, boardID int64, roleID discord.RoleID) (bool, error) {
	res, err := db.Exec(addMentionQuery, channelID, boardID, roleID)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// GetMentionRoles returns all VLIVE mentions for a VLIVE feed.
func (db *DB) GetMentionRoles(channelID discord.ChannelID, boardID int64) ([]discord.RoleID, error) {
	var roleIDs []discord.RoleID
	err := db.Select(&roleIDs, getMentionRolesQuery, channelID, boardID)

	return roleIDs, err
}

// GetMentionsByGuild returns all VLIVE mentions in a guild ID.
func (db *DB) GetMentions(channelID discord.ChannelID, boardID int64) ([]Mention, error) {
	var mentionRoles []Mention
	err := db.Select(&mentionRoles, getMentionsQuery, channelID, boardID)

	return mentionRoles, err
}

// GetMentionsByGuild returns all VLIVE mentions in a guild ID.
func (db *DB) GetMentionsByGuild(guildID discord.GuildID) ([]Mention, error) {
	var mentionRoles []Mention
	err := db.Select(&mentionRoles, getMentionsByGuildQuery, guildID)

	return mentionRoles, err
}

// RemoveMention removes a VLIVE mention.
func (db *DB) RemoveMention(
	channelID discord.ChannelID, boardID int64, roleID discord.RoleID) (bool, error) {
	res, err := db.Exec(removeMentionQuery, channelID, boardID, roleID)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// RemoveMentions removes all VLIVE mentions for a VLIVE feed.
func (db *DB) RemoveMentions(
	channelID discord.ChannelID, boardID int64) (bool, error) {
	res, err := db.Exec(removeMentionsQuery, channelID, boardID)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// ClearGuildMentions removes all VLIVE mentions in a guild ID.
func (db *DB) ClearGuildMentions(guildID discord.GuildID) (int64, error) {
	res, err := db.Exec(clearGuildMentionsQuery, guildID)
	if err != nil {
		return 0, err
	}

	cleared, err := res.RowsAffected()
	return cleared, err
}
