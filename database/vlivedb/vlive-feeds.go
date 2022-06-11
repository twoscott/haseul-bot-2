package vlivedb

import "github.com/diamondburned/arikawa/v3/discord"

// Feed represents a VLIVE feed database entry.
type Feed struct {
	GuildID   discord.GuildID   `db:"guildid"`
	ChannelID discord.ChannelID `db:"channelid"`
	BoardID   int64             `db:"vliveboardid"`
}

const (
	createVLIVEFeedsTableQuery = `
		CREATE TABLE IF NOT EXISTS VLIVEFeeds(
			guildID      INT8 NOT NULL,
			channelID    INT8 NOT NULL,
			boardID INT8 NOT NULL,
			PRIMARY KEY(channelID, boardID),
			FOREIGN KEY(boardID) REFERENCES VLIVEBoards(id)
		)`
	addFeedQuery = `
		INSERT INTO VLIVEFeeds VALUES($1, $2, $3) ON CONFLICT DO NOTHING`
	getFeedsByUserQuery = `
		SELECT * FROM VLIVEFeeds WHERE boardID = $1`
	getFeedsByGuildQuery = `
		SELECT * FROM VLIVEFeeds WHERE guildID = $1`
	getFeedQuery = `
		SELECT * FROM VLIVEFeeds WHERE channelID = $1 AND boardID = $2`
	removeFeedQuery = `
		DELETE FROM VLIVEFeeds WHERE channelID = $1 AND boardID = $2`
	removeFeedsByUserQuery = `
		DELETE FROM VLIVEFeeds WHERE boardID = $1`
	removeFeedsByChannelQuery = `
		DELETE FROM VLIVEFeeds WHERE channelID = $1`
	clearGuildFeedsQuery = `
		DELETE FROM VLIVEFeeds WHERE guildID = $1`
)

// AddFeed adds a new Twitter feed to the database.
func (db *DB) AddFeed(
	guildID discord.GuildID,
	channelID discord.ChannelID,
	boardID int64) (bool, error) {

	res, err := db.Exec(addFeedQuery, guildID, channelID, boardID)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// GetFeedsByUser returns all Twitter feeds set up with
// the given Twitter user ID.
func (db *DB) GetFeedsByUser(boardID int64) ([]Feed, error) {
	var feeds []Feed
	err := db.Select(&feeds, getFeedsByUserQuery, boardID)

	return feeds, err
}

// GetFeedsByGuild returns all Twitter feeds set up in the provided guild ID.
func (db *DB) GetFeedsByGuild(guildID discord.GuildID) ([]Feed, error) {
	var feeds []Feed
	err := db.Select(&feeds, getFeedsByGuildQuery, guildID)

	return feeds, err
}

// GetFeed returns a Twitter feed from the database.
func (db *DB) GetFeed(
	channelID discord.ChannelID, boardID int64) (*Feed, error) {

	var feed Feed
	err := db.Get(&feed, getFeedQuery, channelID, boardID)

	return &feed, err
}

// RemoveFeed removes a Twitter feed from the database.
func (db *DB) RemoveFeed(
	channelID discord.ChannelID, boardID int64) (bool, error) {

	res, err := db.Exec(removeFeedQuery, channelID, boardID)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// RemoveFeedsByUser removes all Twitter feeds for a given Twitter user ID.
func (db *DB) RemoveFeedsByUser(boardID int64) (int64, error) {
	res, err := db.Exec(removeFeedsByUserQuery, boardID)
	if err != nil {
		return 0, err
	}

	deleted, err := res.RowsAffected()
	return deleted, err
}

// RemoveFeedsByChannel removes all Twitter feeds for a given channel ID.
func (db *DB) RemoveFeedsByChannel(channelID discord.ChannelID) (int64, error) {
	res, err := db.Exec(removeFeedsByChannelQuery, channelID)
	if err != nil {
		return 0, err
	}

	deleted, err := res.RowsAffected()
	return deleted, err
}

// ClearGuildFeeds removes all Twitter feeds in a given guild ID.
func (db *DB) ClearGuildFeeds(guildID discord.GuildID) (int64, error) {
	res, err := db.Exec(clearGuildFeedsQuery, guildID)
	if err != nil {
		return 0, err
	}

	cleared, err := res.RowsAffected()
	return cleared, err
}
