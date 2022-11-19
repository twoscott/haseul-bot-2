package logs

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var logsMemberChannelCommand = &router.SubCommand{
	Name:        "channel",
	Description: "Sets the channel for member logs to be posted to",
	Handler: &router.CommandHandler{
		Executor: logsMemberChannelExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.ChannelOption{
			OptionName:  "channel",
			Description: "The channel to log member joins & leaves in",
			Required:    true,
			ChannelTypes: []discord.ChannelType{
				discord.GuildText,
				discord.GuildNews,
			},
		},
	},
}

func logsMemberChannelExec(ctx router.CommandCtx) {
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

	_, err := db.Guilds.SetMemberLogsChannel(ctx.Interaction.GuildID, channel.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while setting member logs channel.")
		return
	}

	ctx.RespondSuccess("Member logs channel set to " + channel.Mention() + ".")
}
