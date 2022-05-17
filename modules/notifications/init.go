package notifications

import (
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
)

var db *database.DB

func Init(rt *router.Router) {
	db = database.GetInstance()

	rt.MustRegisterCommand(notiCommand)

	notiCommand.MustRegisterSubCommand(notiGlobCommand)
	notiGlobCommand.MustRegisterSubCommand(notiGlobAddCommand)
	notiGlobCommand.MustRegisterSubCommand(notiGlobRemoveCommand)
	notiGlobCommand.MustRegisterSubCommand(notiGlobClearCommand)

	notiCommand.MustRegisterSubCommand(notiAddCommand)
	notiCommand.MustRegisterSubCommand(notiRemoveCommand)
	notiCommand.MustRegisterSubCommand(notiClearCommand)

	notiCommand.MustRegisterSubCommand(notiListCommand)

	notiCommand.MustRegisterSubCommand(notiGuildCommand)
	notiGuildCommand.MustRegisterSubCommand(notiGuildMuteCommand)
	notiGuildCommand.MustRegisterSubCommand(notiGuildUnmuteCommand)

	notiCommand.MustRegisterSubCommand(notiChannelCommand)
	notiChannelCommand.MustRegisterSubCommand(notiChannelMuteCommand)
	notiChannelCommand.MustRegisterSubCommand(notiChannelUnmuteCommand)

	notiCommand.MustRegisterSubCommand(notiDndCommand)

	rt.RegisterMessageHandler(checkKeywords)
}
