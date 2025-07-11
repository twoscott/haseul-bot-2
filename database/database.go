package database

import (
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/twoscott/haseul-bot-2/database/commanddb"
	"github.com/twoscott/haseul-bot-2/database/guilddb"
	"github.com/twoscott/haseul-bot-2/database/invitedb"
	"github.com/twoscott/haseul-bot-2/database/lastfmdb"
	"github.com/twoscott/haseul-bot-2/database/levelsdb"
	"github.com/twoscott/haseul-bot-2/database/marriagedb"
	"github.com/twoscott/haseul-bot-2/database/notifdb"
	"github.com/twoscott/haseul-bot-2/database/reminderdb"
	"github.com/twoscott/haseul-bot-2/database/repdb"
	"github.com/twoscott/haseul-bot-2/database/rolesdb"
	"github.com/twoscott/haseul-bot-2/database/youtubedb"
)

type DB struct {
	*sqlx.DB
	Commands      *commanddb.DB
	Guilds        *guilddb.DB
	Invites       *invitedb.DB
	LastFM        *lastfmdb.DB
	Levels        *levelsdb.DB
	Marriages     *marriagedb.DB
	Notifications *notifdb.DB
	Reminders     *reminderdb.DB
	Reps          *repdb.DB
	Roles         *rolesdb.DB
	YouTube       *youtubedb.DB
}

var (
	db   *DB
	once sync.Once
)

func GetInstance() *DB {
	once.Do(func() {
		dbConn := mustGetConnection()

		db = &DB{
			DB:            dbConn,
			Commands:      commanddb.New(dbConn),
			Guilds:        guilddb.New(dbConn),
			Invites:       invitedb.New(dbConn),
			LastFM:        lastfmdb.New(dbConn),
			Marriages:     marriagedb.New(dbConn),
			Notifications: notifdb.New(dbConn),
			Reminders:     reminderdb.New(dbConn),
			Reps:          repdb.New(dbConn),
			Roles:         rolesdb.New(dbConn),
			Levels:        levelsdb.New(dbConn),
			YouTube:       youtubedb.New(dbConn),
		}
	})

	return db
}
