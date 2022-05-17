package twitter

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var twtToggleRetweetsCommand = &router.Command{
	Name: "retweets",
	RequiredPermissions: 0 |
		discord.PermissionManageChannels |
		discord.PermissionManageGuild,
	IncludeAdmin: true,
	UseTyping:    true,
	Run:          twtToggleRetweetsRun,
}

func twtToggleRetweetsRun(ctx router.CommandCtx, args []string) {
	if len(args) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a Twitter user and Discord channel to "+
				"toggle retweets for.",
		)
		return
	}
	if len(args) < 2 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a Discord channel to toggle retweets for.",
		)
		return
	}

	screenName := parseUserIfURL(args[0])
	user, ok := fetchUser(ctx, screenName)
	if !ok {
		return
	}

	channel, ok := parseChannelArg(ctx, args[1])
	if !ok {
		return
	}

	toggled, err := db.Twitter.ToggleFeedRetweets(channel.ID, user.ID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while toggling retweets.",
		)
		return
	}
	if !toggled {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			fmt.Sprintf(
				"I could not find a Twitter feed for @%s in %s.",
				user.ScreenName, channel.Mention(),
			),
		)
		return
	}

	feed, err := db.Twitter.GetFeed(channel.ID, user.ID)
	if err != nil {
		dctools.ReplyWithSuccess(ctx.State, ctx.Msg,
			"Retweets were toggled for this feed.",
		)
	}

	var content string
	if feed.Retweets {
		content = fmt.Sprintf(
			"You will now receive retweets from @%s in %s.",
			user.ScreenName, channel.Mention(),
		)
	} else {
		content = fmt.Sprintf(
			"You will no longer receive retweets from @%s in %s.",
			user.ScreenName, channel.Mention(),
		)
	}

	dctools.ReplyWithSuccess(ctx.State, ctx.Msg, content)
}
