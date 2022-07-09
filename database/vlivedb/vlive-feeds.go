package vlivedb

import "github.com/diamondburned/arikawa/v3/discord"

// Feed represents a VLIVE feed database entry.
type Feed struct {
	BoardID   int64             `db:"boardid"`
	GuildID   discord.GuildID   `db:"guildid"`
	ChannelID discord.ChannelID `db:"channelid"`
	PostTypes PostTypes         `db:"posttypes"`
	Reposts   bool              `db:"reposts"`
}

type PostTypes int16

const (
	AllPosts PostTypes = iota
	VideosOnly
	PostsOnly
)

const (
	createVLIVEFeedsTableQuery = `
		CREATE TABLE IF NOT EXISTS VLIVEFeeds(
			boardID   INT8    NOT NULL,
			guildID   INT8    NOT NULL,
			channelID INT8    NOT NULL,
			postTypes INT2    DEFAULT 0,
			reposts   BOOLEAN DEFAULT TRUE,
			PRIMARY KEY(boardID, channelID),
			FOREIGN KEY(boardID) REFERENCES VLIVEBoards(id)
		)`
	addFeedQuery = `
		INSERT INTO VLIVEFeeds VALUES($1, $2, $3, $4, $5) 
		ON CONFLICT DO NOTHING`
	getFeedsByBoardQuery = `
		SELECT * FROM VLIVEFeeds WHERE boardID = $1`
	getFeedsByGuildQuery = `
		SELECT * FROM VLIVEFeeds WHERE guildID = $1`
	getFeedQuery = `
		SELECT * FROM VLIVEFeeds WHERE channelID = $1 AND boardID = $2`
	removeFeedQuery = `
		DELETE FROM VLIVEFeeds WHERE channelID = $1 AND boardID = $2`
	removeFeedsByBoardQuery = `
		DELETE FROM VLIVEFeeds WHERE boardID = $1`
	removeFeedsByChannelQuery = `
		DELETE FROM VLIVEFeeds WHERE channelID = $1`
	clearGuildFeedsQuery = `
		DELETE FROM VLIVEFeeds WHERE guildID = $1`
)

// AddFeed adds a new VLIVE feed to the database.
func (db *DB) AddFeed(
	boardID int64,
	guildID discord.GuildID,
	channelID discord.ChannelID,
	postTypes PostTypes,
	reposts bool) (bool, error) {

	res, err := db.Exec(
		addFeedQuery, boardID, guildID, channelID, postTypes, reposts,
	)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// GetFeedsByBoard returns all VLIVE feeds set up with
// the given VLIVE board ID.
func (db *DB) GetFeedsByBoard(boardID int64) ([]Feed, error) {
	var feeds []Feed
	err := db.Select(&feeds, getFeedsByBoardQuery, boardID)

	return feeds, err
}

// GetFeedsByGuild returns all VLIVE feeds set up in the provided guild ID.
func (db *DB) GetFeedsByGuild(guildID discord.GuildID) ([]Feed, error) {
	var feeds []Feed
	err := db.Select(&feeds, getFeedsByGuildQuery, guildID)

	return feeds, err
}

// GetFeed returns a VLIVE feed from the database.
func (db *DB) GetFeed(
	channelID discord.ChannelID, boardID int64) (*Feed, error) {

	var feed Feed
	err := db.Get(&feed, getFeedQuery, channelID, boardID)

	return &feed, err
}

// RemoveFeed removes a VLIVE feed from the database.
func (db *DB) RemoveFeed(
	channelID discord.ChannelID, boardID int64) (bool, error) {

	res, err := db.Exec(removeFeedQuery, channelID, boardID)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// RemoveFeedsByBoard removes all VLIVE feeds for a given VLIVE board ID.
func (db *DB) RemoveFeedsByBoard(boardID int64) (int64, error) {
	res, err := db.Exec(removeFeedsByBoardQuery, boardID)
	if err != nil {
		return 0, err
	}

	deleted, err := res.RowsAffected()
	return deleted, err
}

// RemoveFeedsByChannel removes all VLIVE feeds for a given channel ID.
func (db *DB) RemoveFeedsByChannel(channelID discord.ChannelID) (int64, error) {
	res, err := db.Exec(removeFeedsByChannelQuery, channelID)
	if err != nil {
		return 0, err
	}

	deleted, err := res.RowsAffected()
	return deleted, err
}

// ClearGuildFeeds removes all VLIVE feeds in a given guild ID.
func (db *DB) ClearGuildFeeds(guildID discord.GuildID) (int64, error) {
	res, err := db.Exec(clearGuildFeedsQuery, guildID)
	if err != nil {
		return 0, err
	}

	cleared, err := res.RowsAffected()
	return cleared, err
}
