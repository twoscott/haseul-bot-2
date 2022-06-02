package vlive

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var vliveCommand = &router.Command{
	Name:        "vlive",
	Description: "Commands pertaining to VLIVE functionality",
	RequiredPermissions: discord.NewPermissions(
		discord.PermissionManageChannels,
		discord.PermissionManageRoles,
	),
}
