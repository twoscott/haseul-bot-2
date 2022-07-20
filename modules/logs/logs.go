package logs

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var logsCommand = &router.Command{
	Name:        "logs",
	Description: "Commands pertaining to server logs settings",
	RequiredPermissions: discord.NewPermissions(
		discord.PermissionManageChannels,
		discord.PermissionManageRoles,
	),
}
