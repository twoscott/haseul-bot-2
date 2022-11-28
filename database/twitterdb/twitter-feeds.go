package twitterdb

import "github.com/diamondburned/arikawa/v3/discord"

// Feed represents a Twitter feed database entry.
type Feed struct {
	GuildID       discord.GuildID   `db:"guildid"`
	ChannelID     discord.ChannelID `db:"channelid"`
	TwitterUserID int64             `db:"twitteruserid"`
	Replies       bool              `db:"replies"`
	Retweets      bool              `db:"retweets"`
}

const (
	createTwitterFeedsTableQuery = `
		CREATE TABLE IF NOT EXISTS TwitterFeeds(
			guildID       INT8    NOT NULL,
			channelID     INT8    NOT NULL,
			twitterUserID INT8    NOT NULL,
			replies		  BOOLEAN NOT NULL DEFAULT TRUE,
			retweets      BOOLEAN NOT NULL DEFAULT TRUE,
			PRIMARY KEY(channelID, twitterUserID),
			FOREIGN KEY(twitterUserID) REFERENCES TwitterUsers(ID)
		)`

	addFeedQuery = `
		INSERT INTO TwitterFeeds VALUES($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	getFeedsByUserQuery = `
		SELECT * FROM TwitterFeeds WHERE twitterUserID = $1`
	getFeedsByGuildQuery = `
		SELECT * FROM TwitterFeeds WHERE guildID = $1`
	getFeedQuery = `
		SELECT * FROM TwitterFeeds WHERE channelID = $1 AND twitterUserID = $2`
	removeFeedQuery = `
		DELETE FROM TwitterFeeds WHERE channelID = $1 AND twitterUserID = $2`
	removeFeedsByUserQuery = `
		DELETE FROM TwitterFeeds WHERE twitterUserID = $1`
	removeFeedsByChannelQuery = `
		DELETE FROM TwitterFeeds WHERE channelID = $1`
	clearGuildFeedsQuery = `
		DELETE FROM TwitterFeeds WHERE guildID = $1`
	toggleRepliesQuery = `
		UPDATE TwitterFeeds SET replies = NOT replies 
		WHERE channelID = $1 AND twitterUserID = $2`
	toggleRetweetsQuery = `
		UPDATE TwitterFeeds SET retweets = NOT retweets 
		WHERE channelID = $1 AND twitterUserID = $2`
)

// AddFeed adds a new Twitter feed to the database.
func (db *DB) AddFeed(
	guildID discord.GuildID,
	channelID discord.ChannelID,
	twitterUserID int64,
	replies bool,
	retweets bool) (bool, error) {

	res, err := db.Exec(
		addFeedQuery, guildID, channelID, twitterUserID, replies, retweets,
	)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// GetFeedsByUser returns all Twitter feeds set up with
// the given Twitter user ID.
func (db *DB) GetFeedsByUser(twitterUserID int64) ([]Feed, error) {
	var feeds []Feed
	err := db.Select(&feeds, getFeedsByUserQuery, twitterUserID)

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
	channelID discord.ChannelID, twitterUserID int64) (*Feed, error) {

	var feed Feed
	err := db.Get(&feed, getFeedQuery, channelID, twitterUserID)

	return &feed, err
}

// RemoveFeed removes a Twitter feed from the database.
func (db *DB) RemoveFeed(
	channelID discord.ChannelID, twitterUserID int64) (bool, error) {

	res, err := db.Exec(removeFeedQuery, channelID, twitterUserID)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// RemoveFeedsByUser removes all Twitter feeds for a given Twitter user ID.
func (db *DB) RemoveFeedsByUser(twitterUserID int64) (int64, error) {
	res, err := db.Exec(removeFeedsByUserQuery, twitterUserID)
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

// ToggleFeedReplies toggles the feed replies setting for a Twitter feed.
func (db *DB) ToggleFeedReplies(
	channelID discord.ChannelID, twitterUserID int64) (bool, error) {

	res, err := db.Exec(toggleRepliesQuery, channelID, twitterUserID)
	if err != nil {
		return false, err
	}

	toggled, err := res.RowsAffected()
	return toggled > 0, err
}

// ToggleFeedRetweets toggles the feed retweets setting for a Twitter feed.
func (db *DB) ToggleFeedRetweets(
	channelID discord.ChannelID, twitterUserID int64) (bool, error) {

	res, err := db.Exec(toggleRetweetsQuery, channelID, twitterUserID)
	if err != nil {
		return false, err
	}

	toggled, err := res.RowsAffected()
	return toggled > 0, err
}
