package logs

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var logsMessageChannelCommand = &router.SubCommand{
	Name:        "channel",
	Description: "Sets the channel for message logs to be posted to",
	Handler: &router.CommandHandler{
		Executor: logsMessageChannelExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.ChannelOption{
			OptionName:   "channel",
			Description:  "The channel to log deleted and edited messages in",
			Required:     true,
			ChannelTypes: dctools.TextChannelTypes(),
		},
	},
}

func logsMessageChannelExec(ctx router.CommandCtx) {
	snowflake, _ := ctx.Options.Find("channel").SnowflakeValue()
	channelID := discord.ChannelID(snowflake)
	if !channelID.IsValid() {
		ctx.RespondWarning(
			"Malformed Discord channel provided.",
		)
		return
	}

	channel, cerr := ctx.ParseSendableChannel(channelID)
	if cerr != nil {
		ctx.RespondCmdMessage(cerr)
		return
	}

	set, err := db.Guilds.SetMessageLogsChannel(ctx.Interaction.GuildID, channel.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while setting message logs channel.")
		return
	}

	if !set {
		err := fmt.Errorf(
			"message logs channel wasn't updated for %d",
			ctx.Interaction.GuildID,
		)
		log.Println(err)
		ctx.RespondError("Error occurred while setting message logs channel.")
		return
	}

	ctx.RespondSuccess("Message logs channel set to " + channel.Mention() + ".")
}
