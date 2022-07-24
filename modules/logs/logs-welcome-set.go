package logs

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var logsWelcomeSetCommand = &router.SubCommand{
	Name:        "set",
	Description: "Sets the channel for welcome messages to be posted to",
	Handler: &router.CommandHandler{
		Executor: logsWelcomeSetExec,
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

func logsWelcomeSetExec(ctx router.CommandCtx) {
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

	_, err := db.Guilds.SetWelcomeChannel(ctx.Interaction.GuildID, channel.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while setting welcome channel.")
		return
	}

	ctx.RespondSuccess("Welcome channel set to " + channel.Mention() + ".")
}
