package lastfm

import (
	"log"

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
	if key == "" || secret == "" {
		log.Fatalln("No Last.fm API key or secret provided in config file")
	}

	lf = lastfm.New(key, secret)

	rt.AddCommand(lastFmCommand)
	lastFmCommand.AddSubCommand(lastFmCurrentCommand)
	lastFmCommand.AddSubCommand(lastFmDeleteCommand)
	lastFmCommand.AddSubCommand(lastFmRecentCommand)
	lastFmCommand.AddSubCommand(lastFmSetCommand)
	lastFmCommand.AddSubCommand(lastFmYouTubeCommand)

	lastFmCommand.AddSubCommandGroup(lastFmTopCommandGroup)
	lastFmTopCommandGroup.AddSubCommand(lastFmTopAlbumsCommand)
	lastFmTopCommandGroup.AddSubCommand(lastFmTopArtistsCommand)
	lastFmTopCommandGroup.AddSubCommand(lastFmTopTracksCommand)

	// TODO:
	//	/fm collage [type] [dimensions] [timeframe]
}
