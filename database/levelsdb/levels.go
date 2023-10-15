package levelsdb

import (
	"math"

	"github.com/diamondburned/arikawa/v3/discord"
)

const (
	logOffset   = 100_000
	logModifier = 200
	logBase     = 10
)

func getLevelFromXP(xp int64) int {
	return int(math.Log10(float64(xp+logOffset)/logOffset) * logModifier)
}

func getRequiredXPForLevel(level int) int64 {
	return int64(math.Pow10(level/logModifier)*logOffset) - logOffset
}

type Progress struct {
	XP           int64
	NextLevelXP  int64
	CurrentLevel int
}

type UserXP struct {
	UserID discord.UserID `db:"userid"`
	XP     int64          `db:"xp"`
}

func (u UserXP) Level() int {
	return getLevelFromXP(u.XP)
}

func (u UserXP) Progress() Progress {
	level := u.Level()

	baseXP := getRequiredXPForLevel(level)
	nextXP := getRequiredXPForLevel(level + 1)

	xpProgress := u.XP - baseXP

	return Progress{
		XP:           xpProgress,
		NextLevelXP:  nextXP,
		CurrentLevel: level,
	}
}

type GuildUserXP struct {
	UserXP  `db:""`
	GuildID discord.GuildID `db:"guildid"`
}

const (
	createUserXPTableQuery = `
		CREATE TABLE IF NOT EXISTS UserXP(
			guildID INT8 NOT NULL,
			userID  INT8 NOT NULL,
			xp      INT8 NOT NULL DEFAULT 0,
			PRIMARY KEY(guildID, userID)
		)`
	addUserXPQuery = `
		INSERT INTO UserXP VALUES($1, $2, $3)
		ON CONFLICT(guildID, userID) DO UPDATE SET xp = UserXP.xp + $3
		RETURNING xp`
	getUserXPQuery = `
		SELECT xp FROM UserXP WHERE guildID = $1 AND userID = $2`
	getGlobalXPQuery = `
		SELECT SUM(xp) FROM UserXP WHERE userID = $1`
	getTopUsersQuery = `
		SELECT * FROM UserXP WHERE guildID = $1
		ORDER BY xp DESC
		LIMIT $2`
	getTopGlobalUsersQuery = `
		SELECT userID, SUM(xp) AS xp FROM UserXP
		GROUP BY userID 
		ORDER BY xp DESC 
		LIMIT $1`
	getEntriesSizeQuery = `
		SELECT COUNT(userID) FROM UserXP 
		WHERE guildID = $1`
	getGlobalEntriesSizeQuery = `
		SELECT COUNT(DISTINCT userID) FROM UserXP`
)

// AddUserXP XP for a user in a guild.
func (db *DB) AddUserXP(
	guildID discord.GuildID,
	userID discord.UserID,
	xpAmount int64) (xp int64, err error) {

	return xp, db.Get(&xp, addUserXPQuery, guildID, userID, xpAmount)
}

// GetUserXP returns the XP for a user in a guild.
func (db *DB) GetUserXP(
	guildID discord.GuildID, userID discord.UserID) (xp int64, err error) {

	return xp, db.Get(&xp, getUserXPQuery, guildID, userID)
}

// GetUserGlobalXP returns the XP for a user across all guilds.
func (db *DB) GetUserGlobalXP(userID discord.UserID) (xp int64, err error) {

	return xp, db.Get(&xp, getGlobalXPQuery, userID)
}

// GetTopUsers returns the top users in a guild.
func (db *DB) GetTopUsers(
	guildID discord.GuildID, limit int64) (users []GuildUserXP, err error) {

	return users, db.Select(&users, getTopUsersQuery, guildID, limit)
}

// GetTopGlobalUsers returns the top users globally.
func (db *DB) GetTopGlobalUsers(limit int64) (users []UserXP, err error) {

	return users, db.Select(&users, getTopGlobalUsersQuery, limit)
}

// GetEntriesSize returns the number of entries in a guild.
func (db *DB) GetEntriesSize(guildID discord.GuildID) (size int64, err error) {
	return size, db.Get(&size, getEntriesSizeQuery, guildID)
}

// GetEntriesSize returns the number of entries in the table.
func (db *DB) GetGlobalEntriesSize() (size int64, err error) {
	return size, db.Get(&size, getGlobalEntriesSizeQuery)
}
