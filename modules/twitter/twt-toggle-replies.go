package twitter

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var twtToggleRepliesCommand = &router.Command{
	Name:                "replies",
	RequiredPermissions: discord.PermissionManageChannels | discord.PermissionManageGuild,
	IncludeAdmin:        true,
	UseTyping:           true,
	Run:                 twtToggleRepliesRun,
}

func twtToggleRepliesRun(ctx router.CommandCtx, args []string) {
	if len(args) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a Twitter user and Discord channel to "+
				"toggle replies for.",
		)
		return
	}
	if len(args) < 2 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a Discord channel to toggle replies for.",
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

	toggled, err := db.Twitter.ToggleFeedReplies(channel.ID, user.ID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while toggling replies.",
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
			"Replies were toggled for this feed.",
		)
	}

	var content string
	if feed.Replies {
		content = fmt.Sprintf(
			"You will now receive replies from @%s in %s.",
			user.ScreenName, channel.Mention(),
		)
	} else {
		content = fmt.Sprintf(
			"You will no longer receive replies from @%s in %s.",
			user.ScreenName, channel.Mention(),
		)
	}

	dctools.ReplyWithSuccess(ctx.State, ctx.Msg, content)
}
