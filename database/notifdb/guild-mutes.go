package notifdb

import "github.com/diamondburned/arikawa/v3/discord"

const (
	createNotiGuildMutesTableQuery = `
		CREATE TABLE IF NOT EXISTS NotiGuildMutes(
			userID  INT8 NOT NULL,
			guildID INT8 NOT NULL,
			PRIMARY KEY(userID, guildID)
		)`
	addGuildMute = `
		INSERT INTO NotiGuildMutes VALUES($1, $2) ON CONFLICT DO NOTHING`
	removeGuildMute = `
		DELETE FROM NotiGuildMutes WHERE userID = $1 AND channelID = $2`
)

// MuteGuild adds a guild to a user's mute list
func (db *DB) MuteGuild(
	userID discord.UserID, guildID discord.GuildID) (bool, error) {

	res, err := db.Exec(addGuildMute, userID, guildID)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// UnmuteGuild removes a guild from a user's mute list
func (db *DB) UnmuteGuild(
	userID discord.UserID, guildID discord.GuildID) (bool, error) {

	res, err := db.Exec(removeGuildMute, userID, guildID)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}
