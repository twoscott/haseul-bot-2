package message

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/bot/extras/arguments"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
)

var messageEditCommand = &router.SubCommand{
	Name:        "edit",
	Description: "Edits a message to a channel.",
	Handler: &router.CommandHandler{
		Executor: messageEditExec,
		Defer:    false,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:  "message",
			Description: "A link to the message to fetch",
			Required:    true,
		},
		&discord.StringOption{
			OptionName:  "content",
			Description: "The content of the message",
			Required:    true,
			MaxLength:   option.NewInt(2000),
		},
	},
}

func messageEditExec(ctx router.CommandCtx) {
	link := ctx.Options.Find("message").String()
	newContent := ctx.Options.Find("content").String()

	url := arguments.ParseMessageURL(link)
	if url == nil {
		ctx.RespondWarning("Invalid Discord message URL given.")
		return
	}

	var permissions discord.Permissions
	bot, err := ctx.State.Me()
	if err == nil {
		permissions, err = ctx.State.Permissions(url.ChannelID, bot.ID)
	}
	if err == nil && !permissions.Has(discord.PermissionViewChannel) {
		ctx.RespondWarning(
			"I do not have permission to view the message's channel.",
		)
		return
	}

	msg, err := ctx.State.Message(url.ChannelID, url.MessageID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching message data.")
		return
	}

	if msg.Author.ID != bot.ID {
		ctx.RespondWarning("Provided message was not sent by me.")
		return
	}

	if msg.Interaction != nil {
		ctx.RespondWarning(
			"Provided message was sent in response to a command.",
		)
		return
	}

	_, err = ctx.State.EditMessage(url.ChannelID, url.MessageID, newContent)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while editing message.")
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf("Message edited successfully. %s", msg.URL()),
	)
}
