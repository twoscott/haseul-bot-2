package message

import (
	"bytes"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var messageSendCommand = &router.SubCommand{
	Name:        "send",
	Description: "Sends a message to a channel.",
	Handler: &router.CommandHandler{
		Executor: messageSendExec,
		Defer:    true,
	},
	Options: []discord.CommandOptionValue{
		&discord.ChannelOption{
			OptionName:  "channel",
			Description: "The channel to send the message to",
			Required:    true,
			ChannelTypes: []discord.ChannelType{
				discord.GuildText,
				discord.GuildNews,
			},
		},
		&discord.StringOption{
			OptionName:  "content",
			Description: "The content of the message",
			Required:    false,
			MaxLength:   option.NewInt(2000),
		},
		&discord.AttachmentOption{
			OptionName:  "attachment",
			Description: "The attachment to add to the message",
			Required:    false,
		},
	},
}

func messageSendExec(ctx router.CommandCtx) {
	snowflake, _ := ctx.Options.Find("channel").SnowflakeValue()
	channelID := discord.ChannelID(snowflake)
	channel, cerr := ctx.ParseSendableChannel(channelID)
	if cerr != nil {
		ctx.RespondCmdMessage(cerr)
		return
	}

	content := ctx.Options.Find("content").String()

	snowflake, _ = ctx.Options.Find("attachment").SnowflakeValue()
	attachmentID := discord.AttachmentID(snowflake)
	attachment, atchPresent := ctx.Command.Resolved.Attachments[attachmentID]
	if content == "" && !atchPresent {
		ctx.RespondWarning(
			"You must either provide message content or an attachment to send.",
		)
		return
	}

	data := api.SendMessageData{}

	if content != "" {
		data.Content = content
	}
	if atchPresent {
		attachmentData, err := dctools.DownloadAttachment(attachment)
		if err != nil {
			log.Println(err)
			ctx.RespondError("Error occurred downloading attachment.")
			return
		}

		reader := bytes.NewReader(attachmentData)

		data.Files = append(data.Files, sendpart.File{
			Name:   attachment.Filename,
			Reader: reader,
		})
	}

	_, err := ctx.State.SendMessageComplex(channel.ID, data)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while sending the message.")
		return
	}

	ctx.RespondSuccess("Message sent.")
}
