package command

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
)

var commandDeleteCommand = &router.SubCommand{
	Name:        "delete",
	Description: "Deletes a custom server command",
	Handler: &router.CommandHandler{
		Executor:      commandDeleteExec,
		Autocompleter: completeCommandRemoveName,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:   "name",
			Description:  "The name of the command to delete",
			MaxLength:    option.NewInt(commandNameLimit),
			Required:     true,
			Autocomplete: true,
		},
	},
}

func commandDeleteExec(ctx router.CommandCtx) {
	name := ctx.Options.Find("name").String()

	ok, err := db.Commands.Delete(ctx.Interaction.GuildID, name)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while removing the command.")
		return
	}

	if !ok {
		ctx.RespondWarning(
			fmt.Sprintf("A command named %s does not exist.", name),
		)
		return
	}

	ctx.RespondSuccess("Command removed.")
}

func completeCommandRemoveName(ctx router.AutocompleteCtx) {
	name := ctx.Options.Find("name").String()

	choices := commandNameAutocomplete(ctx.Interaction.GuildID, name)
	ctx.RespondChoices(choices)
}
