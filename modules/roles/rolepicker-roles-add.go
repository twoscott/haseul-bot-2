package roles

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var rolePickerRolesAdd = &router.SubCommand{
	Name:        "add",
	Description: "Add a role for users to be able to select",
	Handler: &router.CommandHandler{
		Executor:      rolePickerRolesAddExec,
		Autocompleter: tierNameCompleter,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:   "tier",
			Description:  "The tier to add the role to",
			MaxLength:    option.NewInt(32),
			Required:     true,
			Autocomplete: true,
		},
		&discord.RoleOption{
			OptionName:  "role",
			Description: "The role to add to the tier",
			Required:    true,
		},
		&discord.StringOption{
			OptionName:  "description",
			Description: "The description to show for the role option",
			MaxLength:   option.NewInt(100),
		},
	},
}

func rolePickerRolesAddExec(ctx router.CommandCtx) {
	tierName := ctx.Options.Find("tier").String()
	description := ctx.Options.Find("description").String()

	roleVal, _ := ctx.Options.Find("role").SnowflakeValue()
	roleID := discord.RoleID(roleVal)
	if !roleID.IsValid() {
		ctx.RespondWarning("Invalid role provided.")
		return
	}

	tier, err := db.RolePicker.GetTierByName(ctx.Interaction.GuildID, tierName)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondWarning("This role tier does not exist.")
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching role tiers.")
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

	ok, err := db.RolePicker.AddRole(roleID, tier.ID, description)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while adding the new role.")
		return
	}
	if !ok {
		ctx.RespondWarning(fmt.Sprintf(
			"The role %s is already added to '%s'.",
			roleID.Mention(), tier.Title(),
		))
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf(
			"The role %s has been added to the '%s' tier.",
			roleID.Mention(), tier.Title(),
		),
	)
}
