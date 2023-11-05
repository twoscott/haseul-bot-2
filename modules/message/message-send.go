package message

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

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
			OptionName:  "attachment-1",
			Description: "An attachment to add to the message",
			Required:    false,
		},
		&discord.AttachmentOption{
			OptionName:  "attachment-2",
			Description: "An attachment to add to the message",
			Required:    false,
		},
		&discord.AttachmentOption{
			OptionName:  "attachment-3",
			Description: "An attachment to add to the message",
			Required:    false,
		},
		&discord.AttachmentOption{
			OptionName:  "attachment-4",
			Description: "An attachment to add to the message",
			Required:    false,
		},
		&discord.AttachmentOption{
			OptionName:  "attachment-5",
			Description: "An attachment to add to the message",
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

	attachments := make([]discord.Attachment, 0, 5)
	for i := 1; i <= 5; i++ {
		name := "attachment-" + strconv.Itoa(i)
		snowflake, _ = ctx.Options.Find(name).SnowflakeValue()
		attachmentID := discord.AttachmentID(snowflake)
		attachment, ok := ctx.Command.Resolved.Attachments[attachmentID]
		if ok {
			attachments = append(attachments, attachment)
		}
	}

	if content == "" && len(attachments) == 0 {
		ctx.RespondWarning(
			"You must either provide message content or an attachment to send.",
		)
		return
	}

	files := make([]sendpart.File, 0, len(attachments))
	for _, a := range attachments {
		data, err := dctools.DownloadAttachment(a)
		if err != nil {
			log.Println(err)
			ctx.RespondError("Error occurred downloading attachment.")
			return
		}

		reader := bytes.NewReader(data)
		files = append(files, sendpart.File{
			Name:   a.Filename,
			Reader: reader,
		})
	}

	data := api.SendMessageData{Content: content, Files: files}

	msg, err := ctx.State.SendMessageComplex(channel.ID, data)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while sending the message.")
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf("Message sent. %s", msg.URL()),
	)
}
