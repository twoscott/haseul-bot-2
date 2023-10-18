package logs

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var logsWelcomeChannelCommand = &router.SubCommand{
	Name:        "channel",
	Description: "Sets the channel for welcome messages to be posted to",
	Handler: &router.CommandHandler{
		Executor: logsWelcomeChannelExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.ChannelOption{
			OptionName:  "channel",
			Description: "The channel to welcome new members in",
			Required:    true,
			ChannelTypes: []discord.ChannelType{
				discord.GuildText,
				discord.GuildNews,
			},
		},
	},
}

func logsWelcomeChannelExec(ctx router.CommandCtx) {
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

	set, err := db.Guilds.SetWelcomeChannel(ctx.Interaction.GuildID, channel.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while setting welcome channel.")
		return
	}

	if !set {
		err := fmt.Errorf(
			"welcome channel wasn't updated for %d",
			ctx.Interaction.GuildID,
		)
		log.Println(err)
		ctx.RespondError("Error occurred while setting welcome channel.")
		return
	}

	ctx.RespondSuccess("Welcome channel set to " + channel.Mention() + ".")
}
