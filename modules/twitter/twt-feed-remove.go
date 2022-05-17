package twitter

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var twtFeedRemoveCommand = &router.Command{
	Name:                "remove",
	Aliases:             []string{"rm", "delete", "del"},
	RequiredPermissions: discord.PermissionManageChannels,
	IncludeAdmin:        true,
	UseTyping:           true,
	Run:                 twtFeedRemoveRun,
}

func twtFeedRemoveRun(ctx router.CommandCtx, args []string) {
	if len(args) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a Twitter user and Discord channel to "+
				"stop receiving tweets for.",
		)
		return
	}
	if len(args) < 2 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a Discord channel to stop receiving tweets in.",
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

	_, err := db.Twitter.RemoveMentions(channel.ID, user.ID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while removing roles from the database.",
		)
		return
	}

	ok, err = db.Twitter.RemoveFeed(channel.ID, user.ID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			fmt.Sprintf(
				"Error occurred while removing @%s from the database.",
				user.ScreenName,
			),
		)
		return
	}
	if !ok {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			fmt.Sprintf(
				"%s is not set up to receive tweets from @%s.",
				channel.Mention(), user.ScreenName,
			),
		)
		return
	}

	dctools.ReplyWithSuccess(ctx.State, ctx.Msg,
		fmt.Sprintf(
			"You will no longer receive tweets from @%s in %s.",
			user.ScreenName, channel.Mention(),
		),
	)
}
