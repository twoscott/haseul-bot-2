package vlive

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var vlvFeedAddCommand = &router.Command{
	Name:                "add",
	RequiredPermissions: discord.PermissionManageChannels,
	IncludeAdmin:        true,
	UseTyping:           true,
	Run:                 vlvFeedAddRun,
}

func vlvFeedAddRun(ctx router.CommandCtx, args []string) {
	if len(args) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a VLIVE user and Discord channel "+
				"to set up a VLIVE feed for.",
		)
		return
	}
	if len(args) < 2 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a Discord channel to set up a VLIVE feed for.",
		)
		return
	}

	screenName, isCode := parseCodeIfURL(args[0])
	user, ok := fetchUser(ctx, screenName)
	if !ok {
		return
	}

	channel, ok := parseChannelArg(ctx, args[1])
	if !ok {
		return
	}

	botUser, err := ctx.State.Me()
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			fmt.Sprintf("Error occurred checking my permissions in %s.",
				channel.Mention(),
			),
		)
		return
	}

	botPermissions, err := ctx.State.Permissions(channel.ID, botUser.ID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			fmt.Sprintf("Error occurred checking my permissions in %s.",
				channel.Mention(),
			),
		)
		return
	}

	neededPerms := 0 |
		discord.PermissionViewChannel |
		discord.PermissionSendMessages

	if !botPermissions.Has(neededPerms) {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			fmt.Sprintf("I do not have permission to send messages in %s!",
				channel.Mention(),
			),
		)
		return
	}

	_, err = db.Twitter.GetUser(user.ID)
	if err != nil {
		ok := addUser(&ctx, user)
		if !ok {
			return
		}
	}

	_, err = db.Twitter.GetUserByGuild(ctx.Msg.GuildID, user.ID)
	if err != nil {
		ok := checkGuildTwitterCount(&ctx, user.ID)
		if !ok {
			return
		}
	}

	ok, err = db.Twitter.AddFeed(ctx.Msg.GuildID, channel.ID, user.ID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			fmt.Sprintf(
				"Error occurred while adding @%s to the database.",
				user.ScreenName,
			),
		)
		return
	}
	if !ok {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			fmt.Sprintf(
				"%s is already set up to receive tweets from @%s.",
				channel.Mention(), user.ScreenName,
			),
		)
		return
	}

	dctools.ReplyWithSuccess(ctx.State, ctx.Msg,
		fmt.Sprintf(
			"You will now receive tweets from @%s in %s.",
			user.ScreenName, channel.Mention()),
	)
}

func addUser(ctx *router.CommandCtx, user *twitter.User) bool {
	tweets, resp, err := twt.Timelines.UserTimeline(&twitter.UserTimelineParams{
		UserID:          user.ID,
		Count:           1,
		ExcludeReplies:  twitter.Bool(false),
		IncludeRetweets: twitter.Bool(true),
		TrimUser:        twitter.Bool(true),
	})
	if err != nil {
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Unknown error occurred while trying to fetch tweets.",
		)
		return false
	}
	if resp.StatusCode != http.StatusOK {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			fmt.Sprintf(
				"Error occurred while fetching neccesary data from @%s.",
				user.ScreenName,
			),
		)
		return false
	}

	if len(tweets) < 1 {
		return true
	}

	_, err = db.Twitter.AddUser(user.ID, tweets[0].ID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			fmt.Sprintf(
				"Error occurred while adding @%s to the database.",
				user.ScreenName,
			),
		)
		return false
	}

	return true
}
