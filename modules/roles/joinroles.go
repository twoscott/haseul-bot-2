package roles

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var joinRoles = &router.Command{
	Name:        "join-roles",
	Description: "Commands pertaining to roles assigned to new members",
	RequiredPermissions: discord.NewPermissions(
		discord.PermissionManageRoles,
	),
}
