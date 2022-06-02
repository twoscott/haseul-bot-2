package youtube

import (
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
)

var db *database.DB

func Init(rt *router.Router) {
	db = database.GetInstance()

	rt.AddCommand(youTubeCommand)
	youTubeCommand.AddSubCommand(youTubeSearchCommand)

	youTubeCommand.AddSubCommandGroup(youTubeHistoryCommand)
	youTubeHistoryCommand.AddSubCommand(youTubeHistoryClearCommand)
	youTubeHistoryCommand.AddSubCommand(youTubeHistoryToggleCommand)
}
