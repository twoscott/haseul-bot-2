package roles

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
)

var rolePickerRolesRemove = &router.SubCommand{
	Name:        "remove",
	Description: "Removes a role from a role tier",
	Handler: &router.CommandHandler{
		Executor:      rolePickerRolesRemoveExec,
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
	},
}

func rolePickerRolesRemoveExec(ctx router.CommandCtx) {
	tierName := ctx.Options.Find("tier").String()

	roleVal, _ := ctx.Options.Find("role").SnowflakeValue()
	roleID := discord.RoleID(roleVal)
	if !roleID.IsValid() {
		ctx.RespondWarning("Invalid role provided.")
		return
	}

	tier, err := db.Roles.GetTierByName(ctx.Interaction.GuildID, tierName)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondWarning(
			fmt.Sprintf(
				"The tier '%s' has no roles added to it.", tier.Title(),
			),
		)
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching role tiers.")
		return
	}

	removed, err := db.Roles.RemoveRole(roleID, tier.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while removing role.")
		return
	}
	if !removed {
		ctx.RespondWarning(fmt.Sprintf(
			"The role %s is already removed from '%s'.",
			roleID.Mention(), tier.Title(),
		))
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf(
			"The role %s has been removed from the '%s' tier.",
			roleID.Mention(), tier.Title(),
		),
	)
}
