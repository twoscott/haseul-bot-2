package notification

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var notificationChannelMuteCommand = &router.SubCommand{
	Name:        "mute",
	Description: "Mutes notifications for a channel",
	Handler: &router.CommandHandler{
		Executor:  notificationChannelMuteExec,
		Ephemeral: true,
	},
	Options: []discord.CommandOptionValue{
		&discord.ChannelOption{
			OptionName:  "channel",
			Description: "The channel to mute notifications from",
			ChannelTypes: []discord.ChannelType{
				discord.GuildText,
				discord.GuildNews,
			},
			Required: true,
		},
	},
}

func notificationChannelMuteExec(ctx router.CommandCtx) {
	snowflake, _ := ctx.Options.Find("channel").SnowflakeValue()
	channelID := discord.ChannelID(snowflake)
	if !channelID.IsValid() {
		ctx.RespondWarning("Invalid channel provided.")
		return
	}

	channel, err := ctx.State.Channel(channelID)
	if err != nil {
		log.Println(err)
		ctx.RespondWarning(
			"Invalid Discord channel provided.",
		)
		return
	}
	if channel.GuildID != ctx.Interaction.GuildID {
		ctx.RespondWarning(
			"Channel provided must belong to this server.",
		)
		return
	}

	muted, err := db.Notifications.MuteChannel(
		ctx.Interaction.SenderID(), channelID,
	)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while trying to the mute the channel",
		)
		return
	}

	if muted {
		ctx.RespondSuccess(
			"You will no longer be notified for keywords mentioned in " +
				channelID.Mention() + ".",
		)
	} else {
		ctx.RespondWarning(
			channelID.Mention() + " is already muted.",
		)
	}
}
