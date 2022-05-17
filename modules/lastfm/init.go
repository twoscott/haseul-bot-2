package lastfm

import (
	"github.com/shkh/lastfm-go/lastfm"
	"github.com/twoscott/haseul-bot-2/config"
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
)

var (
	db *database.DB
	lf *lastfm.Api
)

func Init(rt *router.Router) {
	db = database.GetInstance()

	cfg := config.GetInstance()
	key := cfg.LastFm.Key
	secret := cfg.LastFm.Secret
	lf = lastfm.New(key, secret)

	rt.MustRegisterCommand(fmCommand)
	fmCommand.MustRegisterSubCommand(fmDeleteCommand)
	fmCommand.MustRegisterSubCommand(fmNpCommand)
	fmCommand.MustRegisterSubCommand(fmRecentCommand)
	fmCommand.MustRegisterSubCommand(fmSetCommand)
	fmCommand.MustRegisterSubCommand(fmTaCommand)
	fmCommand.MustRegisterSubCommand(fmTalbCommand)
	fmCommand.MustRegisterSubCommand(fmTtCommand)
	fmCommand.MustRegisterSubCommand(fmYtCommand)
}
