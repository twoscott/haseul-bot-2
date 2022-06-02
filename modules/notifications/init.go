package notifications

import (
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
)

var db *database.DB

func Init(rt *router.Router) {
	db = database.GetInstance()

	rt.AddMessageHandler(checkKeywords)

	rt.AddCommand(notiCommand)
	notiCommand.AddSubCommand(notiAddCommand)
	notiCommand.AddSubCommand(notiClearCommand)
	notiCommand.AddSubCommand(notiDndCommand)
	notiCommand.AddSubCommand(notiListCommand)
	notiCommand.AddSubCommand(notiRemoveCommand)

	notiCommand.AddSubCommandGroup(notiChannelCommand)
	notiChannelCommand.AddSubCommand(notiChannelMuteCommand)
	notiChannelCommand.AddSubCommand(notiChannelUnmuteCommand)
}
