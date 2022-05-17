package notifications

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var notiChannelMuteCommand = &router.Command{
	Name:      "mute",
	Aliases:   []string{"blacklist", "ignore"},
	UseTyping: true,
	Run:       notiChannelMuteRun,
}

func notiChannelMuteRun(ctx router.CommandCtx, args []string) {
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

	muted, err := db.Notifications.MuteChannel(
		ctx.Msg.Author.ID, ctx.Msg.ChannelID,
	)
	if err != nil {
		log.Println(err)
		dctools.SendError(ctx.State, ctx.Msg.ChannelID,
			"Error occurred while trying to the mute the channel",
		)
		return
	}

	if muted {
		dctools.SendSuccess(ctx.State, ctx.Msg.ChannelID,
			"You will no longer be notified for keywords mentioned in "+
				channelID.Mention()+".",
		)
	} else {
		dctools.SendWarning(ctx.State, ctx.Msg.ChannelID,
			channelID.Mention()+" is already muted.",
		)
	}
}
