package message

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/bot/extras/arguments"
	"github.com/twoscott/haseul-bot-2/router"
)

var messageFetchCommand = &router.SubCommand{
	Name:        "fetch",
	Description: "Fetches a message to a channel.",
	Handler: &router.CommandHandler{
		Executor: messageFetchExec,
		Defer:    false,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:  "message",
			Description: "A link to the message to fetch",
			Required:    true,
		},
		&discord.BooleanOption{
			OptionName:  "markdown",
			Description: "Display the message with formatting characters",
		},
	},
}

func messageFetchExec(ctx router.CommandCtx) {
	link := ctx.Options.Find("message").String()
	md, _ := ctx.Options.Find("markdown").BoolValue()

	url := arguments.ParseMessageURL(link)
	if url == nil {
		ctx.RespondWarning("Invalid Discord message URL given.")
		return
	}

	bot, err := ctx.State.Me()
	if err == nil {
		permissions, err := ctx.State.Permissions(url.ChannelID, bot.ID)
		if err == nil && !permissions.Has(discord.PermissionViewChannel) {
			ctx.RespondWarning(
				"I do not have permission to view the message's channel.",
			)
			return
		}
	}

	msg, err := ctx.State.Message(url.ChannelID, url.MessageID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching message data.")
		return
	}

	content := msg.Content
	if md {
		content = fmt.Sprintf("```%s```", content)
	}

	ctx.RespondSimple(content)
}