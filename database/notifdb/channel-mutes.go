package notifdb

import "github.com/diamondburned/arikawa/v3/discord"

const (
	createNotiChannelMutesTableQuery = `
		CREATE TABLE IF NOT EXISTS NotiChannelMutes(
			userID    INT8 NOT NULL,
			channelID INT8 NOT NULL,
			PRIMARY KEY(userID, channelID)
		)`
	addChannelMute = `
		INSERT INTO NotiChannelMutes VALUES($1, $2) ON CONFLICT DO NOTHING`
	removeChannelMute = `
		DELETE FROM NotiChannelMutes WHERE userID = $1 AND channelID = $2`
)

// MuteChannel adds a channel to a user's mute list
func (db *DB) MuteChannel(
	userID discord.UserID, channelID discord.ChannelID) (bool, error) {

	res, err := db.Exec(addChannelMute, userID, channelID)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// UnmuteChannel removes a channel from a user's mute list
func (db *DB) UnmuteChannel(
	userID discord.UserID, channelID discord.ChannelID) (bool, error) {

	res, err := db.Exec(removeChannelMute, userID, channelID)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}
