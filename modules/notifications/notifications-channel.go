package notifications

import "github.com/twoscott/haseul-bot-2/router"

var notificationsChannelCommand = &router.SubCommandGroup{
	Name:        "channel",
	Description: "Commands pertaining to notifications in channels",
}
