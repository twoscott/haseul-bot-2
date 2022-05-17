package twitter

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

const (
	twitterIcon   = "https://abs.twimg.com/icons/apple-touch-icon-192x192.png"
	twitterColour = 0x1da1f2
)

var URLRegex = regexp.MustCompile("https://twitter.com/([^/]+)/?.*")

func parseUserIfURL(user string) string {
	if user == "" {
		return user
	}

	match := URLRegex.FindStringSubmatch(user)
	if match == nil {
		if strings.HasPrefix(user, "@") {
			user = user[1:]
		}
		return user
	}

	return match[1]
}

func fetchUser(ctx router.CommandCtx, screenName string) (*twitter.User, bool) {
	user, resp, err := twt.Users.Show(&twitter.UserShowParams{
		ScreenName: screenName,
	})
	switch err.(type) {
	case nil:
		break
	case twitter.APIError:
		switch resp.StatusCode {
		case http.StatusNotFound:
			dctools.ReplyWithWarning(ctx.State, ctx.Msg,
				fmt.Sprintf("I could not find a user named @%s.", screenName),
			)
		case http.StatusForbidden:
			dctools.ReplyWithWarning(ctx.State, ctx.Msg,
				fmt.Sprintf("@%s is either private or suspended.", screenName),
			)
		default:
			log.Println(err)
			dctools.ReplyWithError(ctx.State, ctx.Msg,
				fmt.Sprintf("Unknown error occurred while trying to find @%s.",
					screenName,
				),
			)
		}
		return nil, false
	default:
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			fmt.Sprintf("Unknown error occurred while trying to find @%s.",
				screenName,
			),
		)
		return nil, false
	}

	return user, true
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
