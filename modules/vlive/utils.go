package vlive

import (
	"regexp"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var URLRegex = regexp.MustCompile("https://www.vlive.tv/channel/([^/]+)/?.*")

func parseCodeIfURL(url string) (string, bool) {
	if url == "" {
		return url, false
	}

	match := URLRegex.FindStringSubmatch(url)
	if match == nil {
		return url, false
	}

	return match[1], true
}

func parseChannelArg(
	ctx router.CommandCtx, channelArg string) (*discord.Channel, bool) {

	channelID := dctools.ParseChannelID(channelArg)
	if !channelID.IsValid() {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Malformed Discord channel provided.",
		)
		return nil, false
	}

	channel, err := ctx.State.Channel(channelID)
	if dctools.ErrMissingAccess(err) {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"I cannot access this channel.",
		)
		return nil, false
	}
	if err != nil {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Invalid Discord channel provided.",
		)
		return nil, false
	}
	if channel.GuildID != ctx.Msg.GuildID {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Channel provided must belong to this server.",
		)
		return nil, false
	}
	if !dctools.IsTextChannel(channel.Type) {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Channel provided must be a text channel.",
		)
		return nil, false
	}

	return channel, true
}

func fetchUser(channelCode string) {}
