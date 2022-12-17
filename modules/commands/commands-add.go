package commands

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/haseul-bot-2/router"
)

const commandLimit = 1024

var commandNameRegex = regexp.MustCompile(`^[\p{L}\p{N}]+$`)

var commandsAddCommand = &router.SubCommand{
	Name:        "add",
	Description: "Adds a custom server command",
	Handler: &router.CommandHandler{
		Executor: commandsAddExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:  "name",
			Description: "The name that will be used to trigger the command",
			MaxLength:   option.NewInt(commandNameLimit),
			Required:    true,
		},
		&discord.StringOption{
			OptionName: "content",
			Description: "The content the bot will respond with when the " +
				"command is triggered",
			MaxLength: option.NewInt(1024),
			Required:  true,
		},
	},
}

func commandsAddExec(ctx router.CommandCtx) {
	name := ctx.Options.Find("name").String()
	validName := commandNameRegex.MatchString(name)
	if !validName {
		ctx.RespondWarning(
			"Command names must consist of only letters and numbers.",
		)
		return
	}

	content := ctx.Options.Find("content").String()

	commands, err := db.Commands.GetAllByGuild(ctx.Interaction.GuildID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		ctx.RespondError("Error occurred while checking existing commands.")
		return
	}
	if len(commands) >= commandLimit {
		ctx.RespondWarning(
			fmt.Sprintf(
				"You cannot have more than %s custom commands in a server",
				humanize.Comma(commandLimit),
			),
		)
		return
	}

	ok, err := db.Commands.Add(ctx.Interaction.GuildID, name, content)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while adding the command to the database",
		)
		return
	}
	if !ok {
		ctx.RespondWarning(
			fmt.Sprintf(
				"A command called '%s' already exists in this server.", name,
			),
		)
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf(
			"Command '%s' added. Use %s to trigger the command.",
			name,
			commandsListCommand.Mention(),
		),
	)
}
