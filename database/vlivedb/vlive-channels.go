package vlivedb

import "github.com/diamondburned/arikawa/v3/discord"

// Channel represents a VLIVE channel database entry.
type Channel struct {
	Code              int64 `db:"id"`
	LastPostTimestamp int64 `db:"lastposttimestamp"`
}

const (
	createVLIVEChannelsTableQuery = `
		CREATE TABLE IF NOT EXISTS VLIVEChannels(
			code              TEXT NOT NULL,
			lastPostTimestamp INT8 DEFAULT  0,
			PRIMARY KEY(code)
		)
	`

	addChannelQuery = `
		INSERT INTO VLIVEChannels VALUES($1, $2) ON CONFLICT DO NOTHING`
	setLastTimestampQuery = `
		UPDATE VLIVEChannels SET lastPostTimestamp = $2 WHERE code = $1`

	getChannelQuery     = `SELECT * FROM VLIVEChannels WHERE code = $1`
	getAllChannelsQuery = `SELECT * FROM VLIVEChannels`
	removeChannelQuery  = `DELETE FROM VLIVEChannels WHERE code = $1`

	getGuildChannelCountQuery = `
		SELECT COUNT(*) FROM VLIVEChannels WHERE code IN (
			SELECT vliveChannelCode FROM VLIVEFeeds WHERE guildID = $1
		)
	`
	getChannelByGuildQuery = `
		SELECT * FROM VLIVEChannels WHERE code = $2 AND code IN (
			SELECT vliveChannelCode FROM VLIVEFeeds WHERE guildID = $1
		)
	`
	getChannelsByGuildQuery = `
		SELECT * FROM VLIVEChannels WHERE code IN (
			SELECT vliveChannelCode FROM VLIVEFeeds WHERE guildID = $1
		)
	`
)

// AddChannel adds a VLIVE channel to the database.
func (db *DB) AddChannel(
	vliveChannelCode string, lastTweetID int64) (bool, error) {

	res, err := db.Exec(addChannelQuery, vliveChannelCode, lastTweetID)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// GetChannel returns a VLIVE channel from the database.
func (db *DB) GetChannel(vliveChannelCode string) (*Channel, error) {
	var vliveChannel Channel
	err := db.Get(&vliveChannel, getChannelQuery, vliveChannelCode)

	return &vliveChannel, err
}

// GetAllChannels returns all VLIVE channels from the database.
func (db *DB) GetAllChannels() ([]Channel, error) {
	var vliveChannels []Channel
	err := db.Select(&vliveChannels, getAllChannelsQuery)

	return vliveChannels, err
}

// GetChannelByGuild returns a VLIVE channel if the guild ID has a VLIVE feed
// set up for the provided VLIVE channel ID.
func (db *DB) GetChannelByGuild(
	guildID discord.GuildID, vliveChannelCode string) (Channel, error) {

	var vliveChannel Channel
	err := db.Get(
		&vliveChannel, getChannelByGuildQuery, guildID, vliveChannelCode,
	)

	return vliveChannel, err
}

// GetChannelsByGuild returns all VLIVE channels that are included in one or
// more VLIVE feeds in the given guild ID.
func (db *DB) GetChannelsByGuild(guildID discord.GuildID) ([]Channel, error) {
	var vliveChannels []Channel
	err := db.Select(&vliveChannels, getChannelsByGuildQuery, guildID)

	return vliveChannels, err
}

// GetGuildChannelCount returns how many unique VLIVE channels a guild has
// VLIVE feeds set up for.
func (db *DB) GetGuildChannelCount(guildID discord.GuildID) (int, error) {
	var vliveCount int
	err := db.Get(&vliveCount, getGuildChannelCountQuery, guildID)

	return vliveCount, err
}

// RemoveChannel removes a VLIVE channel from the database.
func (db *DB) RemoveChannel(vliveChannelCode string) (bool, error) {
	res, err := db.Exec(removeChannelQuery, vliveChannelCode)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// SetLastTimestamp updates the last timestamp for a given VLIVE channel code.
func (db *DB) SetLastTimestamp(
	vliveChannelCode string, timestamp int64) (int64, error) {

	res, err := db.Exec(setLastTimestampQuery, vliveChannelCode, timestamp)
	if err != nil {
		return 0, err
	}

	cleared, err := res.RowsAffected()
	return cleared, err
}
