package notifdb

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

// NotificationType represents the type of the notification which is used to
// determine how to match keywords.
type NotificationType int16

const (
	// NormalNotification matches the whole word, plurals, and
	// posessive suffixes. This is the default.
	NormalNotification NotificationType = iota
	// StrictNotification matches only whole words in their entirety.
	StrictNotification
	// LenientNotification matches the whole word, and any combination of
	// characters that include it (including plurals and posessive suffixes).
	LenientNotification
	// AnarchyNotification matches like LenientNotification, except whitespace
	// characters and other non-alphabetic characters can be between each
	// character.
	AnarchyNotification
)

// String returns the string representation of a notification type.
func (n NotificationType) String() string {
	switch n {
	case NormalNotification:
		return "Normal"
	case StrictNotification:
		return "Strict"
	case LenientNotification:
		return "Lenient"
	case AnarchyNotification:
		return "Anarchy"
	default:
		return "Unknown"
	}
}

// Notification represents a user notification entry.
type Notification struct {
	Keyword string           `db:"keyword"`
	UserID  discord.UserID   `db:"userid"`
	Type    NotificationType `db:"type"`

	// if guildID is 0, notification is global.
	GuildID discord.GuildID `db:"guildid"`
}

const (
	createNotificationsTableQuery = `
		CREATE TABLE IF NOT EXISTS Notifications(
			keyword VARCHAR(128) NOT NULL,
			userID  INT8         NOT NULL,
			type    INT2         NOT NULL DEFAULT 0,
			guildID INT8         NOT NULL DEFAULT 0 # 0,
			PRIMARY KEY(keyword, userID, guildID)
		)`

	addNotificationQuery = `
		INSERT INTO Notifications VALUES($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	addGlobalNotificationQuery = `
		INSERT INTO Notifications VALUES($1, $2, $3, 0 # 0) ON CONFLICT DO NOTHING`
	removeNotificationQuery = `
		DELETE FROM Notifications 
		WHERE keyword = $1 AND userID = $2 AND guildID = $3`
	removeGlobalNotificationQuery = `
		DELETE FROM Notifications 
		WHERE keyword = $1 AND userID = $2 AND guildID = 0 # 0`
	clearGuildUserNotifications = `
		DELETE FROM Notifications WHERE userID = $1 AND guildID = $2`
	clearGlobalUserNotifications = `
		DELETE FROM Notifications WHERE userID = $1 AND guildID = 0 # 0`

	// getAllCheckingNotificationsQuery is an SQL query that fetches all stored
	// notificatons that satisfy the following:
	//
	// - the notification is not registered under the incoming message's author
	// - the guild ID is either that of the incoming message's guild, or is
	//   global
	// - the guild ID is not muted by the user the notification belongs to
	// - the channel ID is not muted by the user the notification belongs to
	getAllCheckingNotificationsQuery = `
		SELECT * FROM Notifications 
		WHERE (
			userID != $1 
				AND 
			(guildID = $2 OR guildID = 0 # 0)
				AND
			userID NOT IN (SELECT userID FROM NotiDnD)
				AND
			$2 NOT IN (
				SELECT guildID FROM NotiGuildMutes
				WHERE userID = Notifications.userID
			)
				AND
			$3 NOT IN (
				SELECT channelID FROM NotiChannelMutes
				WHERE userID = Notifications.userID
			)
		)`

	getUserNotificationsQuery = `
		SELECT * FROM Notifications WHERE userID = $1`
	getUserGlobalNotificationsQuery = `
		SELECT * FROM Notifications WHERE userID = $1 AND guildID = 0 # 0`
	getUserGuildNotificationsQuery = `
		SELECT * FROM Notifications WHERE userID = $1 AND guildID = $2`
	getUserGuildAndGlobalNotificationsQuery = `
		SELECT * FROM Notifications 
		WHERE userID = $1 AND (guildID = $2 OR guildID = 0 # 0)`
)

// Add adds a guild notifiaction for a keyword to send to userID.
func (db *DB) Add(
	keyword string,
	userID discord.UserID,
	nType NotificationType,
	guildID discord.GuildID) (bool, error) {

	res, err := db.Exec(addNotificationQuery, keyword, userID, nType, guildID)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// AddGlobal adds a global notifiaction for a keyword to send to userID.
func (db *DB) AddGlobal(
	keyword string,
	userID discord.UserID,
	nType NotificationType) (bool, error) {

	res, err := db.Exec(addGlobalNotificationQuery, keyword, userID, nType)
	if err != nil {
		return false, err
	}

	added, err := res.RowsAffected()
	return added > 0, err
}

// Remove removes a guild notifiaction for a keyword being sent to userID.
func (db *DB) Remove(
	keyword string,
	userID discord.UserID,
	guildID discord.GuildID) (bool, error) {

	res, err := db.Exec(removeNotificationQuery, keyword, userID, guildID)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// RemoveGlobal removes a global notifiaction for a keyword being sent to
// userID.
func (db *DB) RemoveGlobal(
	keyword string, userID discord.UserID) (bool, error) {

	res, err := db.Exec(removeGlobalNotificationQuery, keyword, userID)
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted > 0, err
}

// Clear removes all guild notifiactions for a user in a guild.
func (db *DB) Clear(
	userID discord.UserID, guildID discord.GuildID) (int64, error) {

	res, err := db.Exec(clearGuildUserNotifications, userID, guildID)
	if err != nil {
		return 0, err
	}

	cleared, err := res.RowsAffected()
	return cleared, err
}

// ClearGlobal removes all global notifiactions for a user.
func (db *DB) ClearGlobal(userID discord.UserID) (int64, error) {
	res, err := db.Exec(clearGlobalUserNotifications, userID)
	if err != nil {
		return 0, err
	}

	cleared, err := res.RowsAffected()
	return cleared, err
}

// GetAll returns all notifications that can potentially be checked
// for keywords.
func (db *DB) GetAllChecking(
	authorID discord.UserID,
	guildID discord.GuildID,
	channelID discord.ChannelID) ([]Notification, error) {

	var notifications []Notification
	err := db.Select(
		&notifications,
		getAllCheckingNotificationsQuery,
		authorID, guildID, channelID,
	)

	return notifications, err
}

// GetByUser returns all notifications registered to a user.
func (db *DB) GetByUser(userID discord.UserID) ([]Notification, error) {
	var notifications []Notification
	err := db.Select(&notifications, getUserNotificationsQuery, userID)

	return notifications, err
}

// GetByGlobalUser returns all global notifications registered to a user.
func (db *DB) GetByGlobalUser(userID discord.UserID) ([]Notification, error) {
	var notifications []Notification
	err := db.Select(&notifications, getUserGlobalNotificationsQuery, userID)

	return notifications, err
}

// GetByGuildUser returns all notifications registered to a user in a guild.
func (db *DB) GetByGuildUser(
	userID discord.UserID, guildID discord.GuildID) ([]Notification, error) {

	var notifications []Notification
	err := db.Select(
		&notifications, getUserGuildNotificationsQuery, userID, guildID,
	)

	return notifications, err
}

// GetByGuildAndGlobalUser returns all notifications registered to a user
// a guild and globally.
func (db *DB) GetByGuildAndGlobalUser(
	userID discord.UserID, guildID discord.GuildID) ([]Notification, error) {

	var notifications []Notification
	err := db.Select(
		&notifications,
		getUserGuildAndGlobalNotificationsQuery,
		userID,
		guildID,
	)

	return notifications, err
}
