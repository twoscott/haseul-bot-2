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

	rt.AddCommand(fmCommand)
	fmCommand.AddSubCommand(fmCurrentCommand)
	fmCommand.AddSubCommand(fmDeleteCommand)
	fmCommand.AddSubCommand(fmRecentCommand)
	fmCommand.AddSubCommand(fmSetCommand)
	fmCommand.AddSubCommand(fmYouTubeCommand)

	fmCommand.AddSubCommandGroup(fmTopCommandGroup)
	fmTopCommandGroup.AddSubCommand(fmTopAlbumsCommand)
	fmTopCommandGroup.AddSubCommand(fmTopArtistsCommand)
	fmTopCommandGroup.AddSubCommand(fmTopTracksCommand)

	// TODO:
	//	/fm collage [type] [dimensions] [timeframe]
}
