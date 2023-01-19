package roles

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var joinRolesAdd = &router.SubCommand{
	Name:        "add",
	Description: "Add a role to be assigned to new members who join the server",
	Handler: &router.CommandHandler{
		Executor: joinRolesAddExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.RoleOption{
			OptionName:  "role",
			Description: "The role to be added to new members",
			Required:    true,
		},
	},
}

func joinRolesAddExec(ctx router.CommandCtx) {
	snowflake, _ := ctx.Options.Find("role").SnowflakeValue()
	roleID := discord.RoleID(snowflake)
	if !roleID.IsValid() {
		ctx.RespondWarning("Malformed role ID provided.")
		return
	}

	if dctools.IsEveryoneRole(ctx.Interaction.GuildID, roleID) {
		ctx.RespondWarning("I cannot assign `@everyone` roles!")
		return
	}

	role, err := ctx.State.Role(ctx.Interaction.GuildID, roleID)
	if err != nil {
		ctx.RespondError("Error occurred fetching role.")
		return
	}

	if role.Managed {
		ctx.RespondWarning("I cannot assign managed bot roles to users!")
		return
	}

	botUser, err := ctx.State.Me()
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while checking role permissions.")
		return
	}

	botCanModify, err := dctools.MemberCanModifyRole(
		ctx.State,
		ctx.Interaction.GuildID,
		ctx.Interaction.ChannelID,
		botUser.ID,
		roleID,
	)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while checking role permissions.")
		return
	}
	if !botCanModify {
		ctx.RespondWarning(
			"I cannot assign roles that are positioned above me in " +
				"the role order!",
		)
		return
	}

	senderCanModify, err := dctools.MemberCanModifyRole(
		ctx.State,
		ctx.Interaction.GuildID,
		ctx.Interaction.ChannelID,
		ctx.Interaction.SenderID(),
		roleID,
	)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while checking role permissions.")
		return
	}
	if !senderCanModify {
		ctx.RespondWarning(
			"You cannot add roles that are positioned above you in " +
				"the role order!",
		)
		return
	}

	ok, err := db.Roles.AddJoinRole(ctx.Interaction.GuildID, roleID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while adding join role.")
		return
	}
	if !ok {
		ctx.RespondWarning(
			"This role is already set to be assigned to new members",
		)
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf(
			"%s will now be assigned to new members who join the server.",
			roleID.Mention(),
		),
	)
}
