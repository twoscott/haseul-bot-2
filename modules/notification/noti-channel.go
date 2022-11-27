package notification

import "github.com/twoscott/haseul-bot-2/router"

var notificationChannelCommand = &router.SubCommandGroup{
	Name:        "channel",
	Description: "Commands pertaining to notifications in channels",
}
