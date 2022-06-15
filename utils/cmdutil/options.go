package cmdutil

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

func ParseAccessibleChannel(
	ctx router.CommandCtx,
	channelID discord.ChannelID) (*discord.Channel, router.CmdResponse) {

	if !channelID.IsValid() {
		return nil, router.Warningf("Malformed Discord channel provided.")
	}

	channel, err := ctx.State.Channel(channelID)
	if dctools.ErrMissingAccess(err) {
		return nil, router.Warningf("I cannot access this channel.")
	}
	if err != nil {
		return nil, router.Warningf("Invalid Discord channel provided.")
	}
	if channel.GuildID != ctx.Interaction.GuildID {
		return nil, router.Warningf(
			"Channel provided must belong to this server.",
		)
	}
	if !dctools.IsTextChannel(channel.Type) {
		return nil, router.Warningf("Channel provided must be a text channel.")
	}

	return channel, nil
}
