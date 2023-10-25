package commands

import (
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

const commandNameLimit = 32

var commandsCommand = &router.Command{
	Name:        "commands",
	Description: "Commands pertaining to custom server commands",
	RequiredPermissions: discord.NewPermissions(
		discord.PermissionManageMessages,
	),
}

func commandNameAutocomplete(guildID discord.GuildID, name string) api.AutocompleteChoices {
	commands, err := db.Commands.GetAllByGuild(guildID)
	if err != nil {
		log.Println(err)
		return nil
	}
	if len(commands) == 0 {
		return nil
	}

	names := make([]string, 0, len(commands))
	for _, cmd := range commands {
		names = append(names, cmd.Name)
	}

	var choices api.AutocompleteStringChoices
	if name == "" {
		choices = dctools.MakeStringChoices(names)
	} else {
		matches := util.SearchSort(names, name)
		choices = dctools.MakeStringChoices(matches)
	}

	return choices
}
