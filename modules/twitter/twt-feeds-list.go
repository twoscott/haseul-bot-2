package twitter

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
	"golang.org/x/exp/slices"
)

var twtFeedsListCommand = &router.SubCommand{
	Name:        "list",
	Description: "Lists all Twitter feeds added to the server",
	Handler: &router.CommandHandler{
		Executor: twtFeedListExec,
		Defer:    true,
	},
}

func twtFeedListExec(ctx router.CommandCtx) {
	feeds, err := db.Twitter.GetFeedsByGuild(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while fetching Twitter feeds from the database.",
		)
		return
	}

	if len(feeds) < 1 {
		ctx.RespondWarning(
			"This server has no Twitter feeds set up.",
		)
		return
	}

	userIDs := make([]int64, 0)
	for _, f := range feeds {
		userIDs = append(userIDs, f.TwitterUserID)
	}

	slices.Sort(userIDs)
	users, resp, err := twt.Users.Lookup(&twitter.UserLookupParams{
		UserID: slices.Compact(userIDs),
	})
	if err != nil || resp.StatusCode != http.StatusOK {
		ctx.RespondError(
			"Unknown error occurred while trying to fetch users.",
		)
		return
	}

	screenNames := make(map[int64]string, len(users))
	for _, u := range users {
		screenNames[u.ID] = u.ScreenName
	}

	feedList := make([]string, 0, len(feeds))
	for _, feed := range feeds {
		feedEntry := feed.ChannelID.Mention() + " - "

		screenName := screenNames[feed.TwitterUserID]
		if screenName == "" {
			feedEntry += fmt.Sprintf("(ID#%d)", feed.TwitterUserID)
		} else {
			handle := fmt.Sprintf("@%s", screenName)
			url := buildUserURL(screenName)
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
	pages := make([]router.MessagePage, len(descriptionPages))
	footer := util.PluraliseWithCount("Feed", int64(len(feeds)))

	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		pages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Name: "Twitter Notifications", Icon: twitterIcon,
					},
					Description: description,
					Color:       twitterColour,
					Footer: &discord.EmbedFooter{
						Text: dctools.SeparateEmbedFooter(
							pageID,
							footer,
						),
					},
				},
			},
		}
	}

	ctx.RespondPaging(pages)
}
