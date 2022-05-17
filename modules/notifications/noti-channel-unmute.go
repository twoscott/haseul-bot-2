package notifications

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var notiChannelUnmuteCommand = &router.Command{
	Name:      "unmute",
	Aliases:   []string{"whitelist", "unignore"},
	UseTyping: true,
	Run:       notiChannelUnmuteRun,
}

func notiChannelUnmuteRun(ctx router.CommandCtx, args []string) {
	var channelID discord.ChannelID

	if len(args) < 1 {
		channelID = ctx.Msg.ChannelID
	} else {
		channelID = dctools.ParseChannelID(args[0])
		if !channelID.IsValid() {
			dctools.ReplyWithWarning(ctx.State, ctx.Msg,
				"Malformed Discord channel provided.",
			)
			return
		}

		channel, err := ctx.State.Channel(channelID)
		if err != nil {
			dctools.ReplyWithWarning(ctx.State, ctx.Msg,
				"Invalid Discord channel provided.",
			)
			return
		}
		if channel.GuildID != ctx.Msg.GuildID {
			dctools.ReplyWithWarning(ctx.State, ctx.Msg,
				"Channel provided must belong to this server.",
			)
			return
		}
		if !dctools.IsTextChannel(channel.Type) {
			dctools.ReplyWithWarning(ctx.State, ctx.Msg,
				"Channel provided must be a text channel.",
			)
			return
		}
	}

	unmuted, err := db.Notifications.UnmuteChannel(
		ctx.Msg.Author.ID, ctx.Msg.ChannelID,
	)
	if err != nil {
		log.Println(err)
		dctools.SendError(ctx.State, ctx.Msg.ChannelID,
			"Error occurred while trying to the mute the channel",
		)
		return
	}

	if unmuted {
		dctools.SendSuccess(ctx.State, ctx.Msg.ChannelID,
			"You will now be notified for keywords mentioned in "+
				channelID.Mention()+".",
		)
	} else {
		dctools.SendWarning(ctx.State, ctx.Msg.ChannelID,
			channelID.Mention()+" is already unmuted.",
		)
	}
}
