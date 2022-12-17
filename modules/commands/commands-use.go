package commands

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
)

var commandsUseCommand = &router.SubCommand{
	Name:        "use",
	Description: "Runs a custom server command",
	Handler: &router.CommandHandler{
		Executor:      commandsUseExec,
		Autocompleter: completeCommandUseName,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:   "name",
			Description:  "The name of the command to use",
			MaxLength:    option.NewInt(commandNameLimit),
			Required:     true,
			Autocomplete: true,
		},
	},
}

func commandsUseExec(ctx router.CommandCtx) {
	name := ctx.Options.Find("name").String()
	content, err := db.Commands.GetContent(ctx.Interaction.GuildID, name)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondWarning(
			fmt.Sprintf("'%s' does not exist.", name),
		)
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching command.")
		return
	}

	err = ctx.RespondText(content)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = db.Commands.Use(ctx.Interaction.GuildID, name)
	if err != nil {
		log.Println(err)
	}
}

func completeCommandUseName(ctx router.AutocompleteCtx) {
	name := ctx.Options.Find("name").String()

	choices := commandNameAutocomplete(ctx.Interaction.GuildID, name)

	ctx.RespondChoices(choices)
}
