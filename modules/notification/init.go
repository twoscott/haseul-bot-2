package notification

import (
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
)

var db *database.DB

func Init(rt *router.Router) {
	db = database.GetInstance()

	rt.AddMessageHandler(checkKeywords)

	rt.AddCommand(notificationCommand)
	notificationCommand.AddSubCommand(notificationAddCommand)
	notificationCommand.AddSubCommand(notificationClearCommand)
	notificationCommand.AddSubCommand(notificationDndCommand)
	notificationCommand.AddSubCommand(notificationListCommand)
	notificationCommand.AddSubCommand(notificationRemoveCommand)

	notificationCommand.AddSubCommandGroup(notificationChannelCommand)
	notificationChannelCommand.AddSubCommand(notificationChannelMuteCommand)
	notificationChannelCommand.AddSubCommand(notificationChannelUnmuteCommand)
}
