package twitter

import (
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var twtFeedListCommand = &router.Command{
	Name:                "list",
	RequiredPermissions: discord.PermissionManageChannels,
	IncludeAdmin:        true,
	UseTyping:           true,
	Run:                 twtFeedListRun,
}

func twtFeedListRun(ctx router.CommandCtx, _ []string) {
	dbUsers, err := db.Twitter.GetUsersByGuild(ctx.Msg.GuildID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while fetching Twitter users from the database.",
		)
		return
	}

	feeds, err := db.Twitter.GetFeedsByGuild(ctx.Msg.GuildID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while fetching Twitter feeds from the database.",
		)
		return
	}

	if len(feeds) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"This server has no Twitter feeds set up.",
		)
		return
	}

	screenNames := make(map[int64]string)
	for _, dbUser := range dbUsers {
		user, _, _ := twt.Users.Show(&twitter.UserShowParams{UserID: dbUser.ID})
		if user != nil {
			screenNames[dbUser.ID] = user.ScreenName
		}
	}

	feedList := make([]string, 0, len(feeds))
	for _, feed := range feeds {
		feedEntry := feed.ChannelID.Mention() + " - "

		screenName := screenNames[feed.TwitterUserID]
		if screenName == "" {
			feedEntry += fmt.Sprintf("(ID#%d)", feed.TwitterUserID)
		} else {
			handle := fmt.Sprintf("@%s", screenName)
			url := fmt.Sprintf("https://twitter.com/%s", screenName)
			feedEntry += dctools.Hyperlink(handle, url)
		}

		if feed.Replies && feed.Retweets {
			feedEntry += " (💬 + ♻️)"
		} else if feed.Replies {
			feedEntry += " (💬)"
		} else if feed.Retweets {
			feedEntry += " (♻️)"
		}

		feedList = append(feedList, feedEntry)
	}

	descriptionPages := util.PagedLines(feedList, 2048, 20)
	messagePages := make([]router.MessagePage, len(descriptionPages))
	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		messagePages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Name: "Twitter Notifications", Icon: twitterIcon,
					},
					Description: description,
					Color:       twitterColour,
					Footer:      &discord.EmbedFooter{Text: pageID},
				},
			},
		}
	}

	cmdutil.ReplyWithPaging(ctx, ctx.Msg, messagePages)
}
