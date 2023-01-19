package roles

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var rolePickerTiersRemove = &router.SubCommand{
	Name:        "remove",
	Description: "Removes a tier from the role picker",
	Handler: &router.CommandHandler{
		Executor:      rolePickerTiersRemoveExec,
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
	},
}

func rolePickerTiersRemoveExec(ctx router.CommandCtx) {
	tierName := ctx.Options.Find("tier").String()

	removed, err := db.Roles.RemoveTier(ctx.Interaction.GuildID, tierName)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while removing tier.")
		return
	}

	formattedName := util.TitleCase(tierName)
	if !removed {
		ctx.RespondWarning(fmt.Sprintf(
			"The tier '%s' is already removed from the role picker.",
			formattedName,
		))
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf(
			"The tier '%s' has been removed from the role picker.",
			formattedName,
		),
	)
}
