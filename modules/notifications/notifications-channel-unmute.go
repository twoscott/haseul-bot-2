package notifications

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var notificationsChannelUnmuteCommand = &router.SubCommand{
	Name:        "unmute",
	Description: "Unmutes notifications for a channel",
	Handler: &router.CommandHandler{
		Executor:  notificationsChannelUnmuteExec,
		Ephemeral: true,
	},
	Options: []discord.CommandOptionValue{
		&discord.ChannelOption{
			OptionName:  "channel",
			Description: "The channel to unmute notifications in",
			ChannelTypes: []discord.ChannelType{
				discord.GuildText,
				discord.GuildNews,
			},
			Required: true,
		},
	},
}

func notificationsChannelUnmuteExec(ctx router.CommandCtx) {
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

	unmuted, err := db.Notifications.UnmuteChannel(
		ctx.Interaction.SenderID(), channelID,
	)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while trying to the unmute the channel",
		)
		return
	}

	if unmuted {
		ctx.RespondSuccess(
			"You will now be notified for keywords mentioned in " +
				channelID.Mention() + ".",
		)
	} else {
		ctx.RespondWarning(
			channelID.Mention() + " is already unmuted.",
		)
	}
}
