package commands

import (
	"log"
	"slices"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/database/commanddb"
	"github.com/twoscott/haseul-bot-2/router"
)

var commandsSearchCommand = &router.SubCommand{
	Name: "search",
	Description: "Searches for a custom command with a given name in " +
		"the server.",
	Handler: &router.CommandHandler{
		Executor:      commandsSearchExec,
		Autocompleter: completeCommandSearchName,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:   "query",
			Description:  "The query to search commands for.",
			MaxLength:    option.NewInt(commandNameLimit),
			Required:     true,
			Autocomplete: true,
		},
	},
}

func commandsSearchExec(ctx router.CommandCtx) {
	query := ctx.Options.Find("query").String()

	commands, err := db.Commands.GetAllByGuild(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching custom commands.")
		return
	}
	if len(commands) == 0 {
		ctx.RespondWarning("This server has no custom commands.")
		return
	}

	filteredCommands := make([]commanddb.Command, 0)
	for _, cmd := range commands {
		if strings.Contains(cmd.Name, query) {
			filteredCommands = append(filteredCommands, cmd)
		}
	}

	slices.SortStableFunc(filteredCommands, func(a, b commanddb.Command) int {
		return len(a.Name) - len(b.Name)
	})
	slices.SortStableFunc(filteredCommands, func(a, b commanddb.Command) int {
		return strings.Index(a.Name, query) - strings.Index(b.Name, query)
	})

	pages := getCommandsListPages(filteredCommands)
	ctx.RespondPaging(pages)
}

func completeCommandSearchName(ctx router.AutocompleteCtx) {
	query := ctx.Options.Find("query").String()

	choices := commandNameAutocomplete(ctx.Interaction.GuildID, query)
	ctx.RespondChoices(choices)
}
