package roles

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var rolePicker = &router.Command{
	Name:        "role-picker",
	Description: "Commands pertaining to the role picker",
	RequiredPermissions: discord.NewPermissions(
		discord.PermissionManageRoles,
	),
}


