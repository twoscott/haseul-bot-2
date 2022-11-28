package twitterdb

import "github.com/diamondburned/arikawa/v3/discord"

// User represents a Twitter user database entry.
type User struct {
	ID          int64 `db:"id"`
	LastTweetID int64 `db:"lasttweetid"`
}

const (
	createTwitterUsersTableQuery = `
		CREATE TABLE IF NOT EXISTS TwitterUsers(
			ID            INT8 NOT NULL,
			LastTweetID   INT8 NOT NULL DEFAULT  0,
			PRIMARY KEY(ID)
		)
	`

	addUserQuery = `
		INSERT INTO TwitterUsers VALUES($1, $2) ON CONFLICT DO NOTHING`
	setLastTweetQuery = `
		UPDATE TwitterUsers SET LastTweetID = $2 WHERE ID = $1`

	getUserQuery     = `SELECT * FROM TwitterUsers WHERE ID = $1`
	getAllUsersQuery = `SELECT * FROM TwitterUsers`
	removeUserQuery  = `DELETE FROM TwitterUsers WHERE ID = $1`

	getGuildUserCountQuery = `
		SELECT COUNT(*) FROM TwitterUsers WHERE ID IN (
			SELECT twitterUserID FROM TwitterFeeds WHERE guildID = $1
		)
	`
	getUserByGuildQuery = `
		SELECT * FROM TwitterUsers WHERE ID = $2 AND ID IN (
			SELECT twitterUserID FROM TwitterFeeds WHERE guildID = $1
		)
	`
	getUsersByGuildQuery = `
		SELECT * FROM TwitterUsers WHERE ID IN (
			SELECT twitterUserID FROM TwitterFeeds WHERE guildID = $1
		)
	`
)

// AddUser adds a Twitter user to the database.
func (db *DB) AddUser(twitterUserID int64, lastTweetID int64) (bool, error) {
	res, err := db.Exec(addUserQuery, twitterUserID, lastTweetID)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// GetUser returns a Twitter user from the database.
func (db *DB) GetUser(twitterUserID int64) (*User, error) {
	var twitterUser User
	err := db.Get(&twitterUser, getUserQuery, twitterUserID)

	return &twitterUser, err
}

// GetAllUsers returns all Twitter users from the database.
func (db *DB) GetAllUsers() ([]User, error) {
	var twitterUsers []User
	err := db.Select(&twitterUsers, getAllUsersQuery)

	return twitterUsers, err
}

// GetUserByGuild returns a Twitter user if the guild ID has a Twitter feed
// set up for the provided Twitter user ID.
func (db *DB) GetUserByGuild(guildID discord.GuildID, twitterUserID int64) (User, error) {
	var twitterUser User
	err := db.Get(&twitterUser, getUserByGuildQuery, guildID, twitterUserID)

	return twitterUser, err
}

// GetUsersByGuild returns all Twitter users that are included in one or more
// Twitter feeds in the given guild ID.
func (db *DB) GetUsersByGuild(guildID discord.GuildID) ([]User, error) {
	var twitterUsers []User
	err := db.Select(&twitterUsers, getUsersByGuildQuery, guildID)

	return twitterUsers, err
}

// GetGuildUserCount returns how many unique Twitter users a guild has
// Twitter feeds set up for.
func (db *DB) GetGuildUserCount(guildID discord.GuildID) (uint64, error) {
	var twitterCount uint64
	err := db.Get(&twitterCount, getGuildUserCountQuery, guildID)

	return twitterCount, err
}

// RemoveUser removes a Twitter user from the database.
func (db *DB) RemoveUser(twitterUserID int64) (bool, error) {
	res, err := db.Exec(removeUserQuery, twitterUserID)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// SetLastTweet updates the last tweet ID for a given twitter user ID.
func (db *DB) SetLastTweet(twitterUserID int64, tweetID int64) (int64, error) {
	res, err := db.Exec(setLastTweetQuery, twitterUserID, tweetID)
	if err != nil {
		return 0, err
	}

	cleared, err := res.RowsAffected()
	return cleared, err
}
