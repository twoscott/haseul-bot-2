package roles

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var joinRolesRemove = &router.SubCommand{
	Name: "remove",
	Description: "Remove a role from being assigned to new members who " +
		"join the server",
	Handler: &router.CommandHandler{
		Executor: joinRolesRemoveExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.RoleOption{
			OptionName:  "role",
			Description: "The role being added to new members",
			Required:    true,
		},
	},
}

func joinRolesRemoveExec(ctx router.CommandCtx) {
	snowflake, _ := ctx.Options.Find("role").SnowflakeValue()
	roleID := discord.RoleID(snowflake)
	if !roleID.IsValid() {
		ctx.RespondWarning("Malformed role ID provided.")
		return
	}

	ok, err := db.Roles.RemoveJoinRole(ctx.Interaction.GuildID, roleID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while removing join role.")
		return
	}
	if !ok {
		ctx.RespondWarning(
			"This role is not set to be assigned to new members",
		)
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf(
			"%s will no longer be assigned to new members who join the server.",
			roleID.Mention(),
		),
	)
}
