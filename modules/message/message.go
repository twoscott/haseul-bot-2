package message

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var messageCommand = &router.Command{
	Name:        "message",
	Description: "Commands pertaining to messages",
	RequiredPermissions: discord.NewPermissions(
		discord.PermissionManageMessages,
	),
}
