package command

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var commandInfoCommand = &router.SubCommand{
	Name:        "info",
	Description: "Displays information about a custom server command",
	Handler: &router.CommandHandler{
		Executor:      commandInfoExec,
		Autocompleter: completeCommandInfoName,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:   "name",
			Description:  "The name of the command to display",
			MaxLength:    option.NewInt(commandNameLimit),
			Required:     true,
			Autocomplete: true,
		},
	},
}

func commandInfoExec(ctx router.CommandCtx) {
	name := ctx.Options.Find("name").String()

	command, err := db.Commands.GetCommand(ctx.Interaction.GuildID, name)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching command data.")
		return
	}

	log.Println(command.Created)

	embed := discord.Embed{
		Author:      &discord.EmbedAuthor{Name: command.Name},
		Description: command.Content,
		Footer: &discord.EmbedFooter{
			Text: util.PluraliseWithCount("Use", command.Uses),
		},
		Color:     dctools.EmbedBackColour,
		Timestamp: discord.Timestamp(command.Created),
	}

	ctx.RespondEmbed(embed)
}

func completeCommandInfoName(ctx router.AutocompleteCtx) {
	name := ctx.Options.Find("name").String()

	choices := commandNameAutocomplete(ctx.Interaction.GuildID, name)
	ctx.RespondChoices(choices)
}
