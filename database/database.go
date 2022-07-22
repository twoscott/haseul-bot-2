package database

import (
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/twoscott/haseul-bot-2/database/guilddb"
	"github.com/twoscott/haseul-bot-2/database/invitedb"
	"github.com/twoscott/haseul-bot-2/database/lastfmdb"
	"github.com/twoscott/haseul-bot-2/database/notifdb"
	"github.com/twoscott/haseul-bot-2/database/twitterdb"
	"github.com/twoscott/haseul-bot-2/database/vlivedb"
	"github.com/twoscott/haseul-bot-2/database/youtubedb"
	"github.com/twoscott/haseul-bot-2/router"
)

type DB struct {
	*sqlx.DB
	LastFM        *lastfmdb.DB
	Invites       *invitedb.DB
	Twitter       *twitterdb.DB
	Guilds        *guilddb.DB
	Notifications *notifdb.DB
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
			Guilds:        guilddb.New(dbConn),
			Invites:       invitedb.New(dbConn),
			LastFM:        lastfmdb.New(dbConn),
			Twitter:       twitterdb.New(dbConn),
			Notifications: notifdb.New(dbConn),
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
