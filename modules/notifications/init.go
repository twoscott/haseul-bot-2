package notifications

import (
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
)

var db *database.DB

func Init(rt *router.Router) {
	db = database.GetInstance()

	rt.AddMessageHandler(checkKeywords)

	rt.AddCommand(notificationsCommand)
	notificationsCommand.AddSubCommand(notificationsAddCommand)
	notificationsCommand.AddSubCommand(notificationsClearCommand)
	notificationsCommand.AddSubCommand(notificationsDndCommand)
	notificationsCommand.AddSubCommand(notificationsListCommand)
	notificationsCommand.AddSubCommand(notificationsDeleteCommand)

	notificationsCommand.AddSubCommandGroup(notificationsChannelCommand)
	notificationsChannelCommand.AddSubCommand(notificationsChannelMuteCommand)
	notificationsChannelCommand.AddSubCommand(notificationsChannelUnmuteCommand)
}
