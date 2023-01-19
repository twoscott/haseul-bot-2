package database

import (
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/twoscott/haseul-bot-2/database/commanddb"
	"github.com/twoscott/haseul-bot-2/database/guilddb"
	"github.com/twoscott/haseul-bot-2/database/invitedb"
	"github.com/twoscott/haseul-bot-2/database/lastfmdb"
	"github.com/twoscott/haseul-bot-2/database/notifdb"
	"github.com/twoscott/haseul-bot-2/database/reminderdb"
	"github.com/twoscott/haseul-bot-2/database/rolesdb"
	"github.com/twoscott/haseul-bot-2/database/twitterdb"
	"github.com/twoscott/haseul-bot-2/database/vlivedb"
	"github.com/twoscott/haseul-bot-2/database/youtubedb"
	"github.com/twoscott/haseul-bot-2/router"
)

type DB struct {
	*sqlx.DB
	Commands      *commanddb.DB
	Guilds        *guilddb.DB
	Invites       *invitedb.DB
	LastFM        *lastfmdb.DB
	Notifications *notifdb.DB
	Reminders     *reminderdb.DB
	Roles         *rolesdb.DB
	Twitter       *twitterdb.DB
	VLIVE         *vlivedb.DB
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
			Notifications: notifdb.New(dbConn),
			Reminders:     reminderdb.New(dbConn),
			Roles:         rolesdb.New(dbConn),
			Twitter:       twitterdb.New(dbConn),
			VLIVE:         vlivedb.New(dbConn),
			YouTube:       youtubedb.New(dbConn),
		}
	})

	return db
}

func Init(rt *router.Router) {
	db := GetInstance()
	db.Guilds.Init(rt)
}
