package lastfm

import (
	"log"

	"github.com/twoscott/gobble-fm/api"
	"github.com/twoscott/haseul-bot-2/config"
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
)

var (
	db *database.DB
	fm *api.Client
)

func Init(rt *router.Router) {
	db = database.GetInstance()
	cfg := config.GetInstance()

	apiKey := cfg.LastFM.Key
	if apiKey == "" {
		log.Fatalln("No Last.fm API key provided in config file")
	}

	fm = api.NewClientKeyOnly(apiKey)

	rt.AddCommand(lastFMCommand)
	lastFMCommand.AddSubCommand(lastFMCollageCommand)
	lastFMCommand.AddSubCommand(lastFMCurrentCommand)
	lastFMCommand.AddSubCommand(lastFMDeleteCommand)
	lastFMCommand.AddSubCommand(lastFMRecentCommand)
	lastFMCommand.AddSubCommand(lastFMSetCommand)
	lastFMCommand.AddSubCommand(lastFMYouTubeCommand)

	lastFMCommand.AddSubCommandGroup(lastFMTopCommandGroup)
	lastFMTopCommandGroup.AddSubCommand(lastFMTopAlbumsCommand)
	lastFMTopCommandGroup.AddSubCommand(lastFMTopArtistsCommand)
	lastFMTopCommandGroup.AddSubCommand(lastFMTopTracksCommand)
}
