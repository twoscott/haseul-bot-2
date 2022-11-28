package command

import (
	"fmt"
	"log"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/database/commandsdb"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
	"golang.org/x/exp/slices"
)

type commandListSort int

const (
	nameSort commandListSort = iota
	usesSort
)

var commandListCommand = &router.SubCommand{
	Name:        "list",
	Description: "Lists all custom commands added to the server",
	Handler: &router.CommandHandler{
		Executor: commandListExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.IntegerOption{
			OptionName:  "sort-by",
			Description: "How to sort the commands",
			Required:    false,
			Choices: []discord.IntegerChoice{
				{Name: "Name", Value: int(nameSort)},
				{Name: "Uses", Value: int(usesSort)},
			},
		},
	},
}

func commandListExec(ctx router.CommandCtx) {
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

	sort, _ := ctx.Options.Find("sort-by").IntValue()
	sortType := commandListSort(sort)

	pages := getCommandsListPages(commands, sortType)
	ctx.RespondPaging(pages)
}

func getCommandsListPages(
	commands []commandsdb.Command, sort commandListSort) []router.MessagePage {

	slices.SortFunc(commands, func(a, b commandsdb.Command) bool {
		return strings.Compare(a.Name, b.Name) < 0
	})

	if sort == usesSort {
		slices.SortStableFunc(commands, func(a, b commandsdb.Command) bool {
			return a.Uses-b.Uses > 0
		})
	}

	commandList := make([]string, 0, len(commands))
	for _, cmd := range commands {
		row := fmt.Sprintf(
			"`%s` (%s)",
			cmd.Name,
			util.PluraliseWithCount("use", cmd.Uses),
		)
		commandList = append(commandList, row)
	}

	descriptionPages := util.PagedLines(commandList, 2048, 20)
	pages := make([]router.MessagePage, len(descriptionPages))
	footer := util.PluraliseWithCount("Command", int64(len(commands)))

	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		pages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Name: "Custom Server Commands",
					},
					Description: description,
					Color:       dctools.EmbedBackColour,
					Footer: &discord.EmbedFooter{
						Text: dctools.SeparateEmbedFooter(
							pageID,
							footer,
						),
					},
				},
			},
		}
	}

	return pages
}
