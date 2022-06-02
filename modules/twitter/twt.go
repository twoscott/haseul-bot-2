package twitter

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var twtCommand = &router.Command{
	Name:        "twitter",
	Description: "Commands pertaining to Twitter functionality",
	RequiredPermissions: discord.NewPermissions(
		discord.PermissionManageChannels,
		discord.PermissionManageRoles,
	),
}
