package vlive

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/database/vlivedb"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
	"github.com/twoscott/haseul-bot-2/utils/vliveutil"
	"golang.org/x/exp/slices"
)

var vliveFeedsListCommand = &router.SubCommand{
	Name:        "list",
	Description: "Lists all VLIVE feeds added to the server",
	Handler: &router.CommandHandler{
		Executor: vliveFeedListExec,
		Defer:    true,
	},
}

func vliveFeedListExec(ctx router.CommandCtx) {
	feeds, err := db.VLIVE.GetFeedsByGuild(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while fetching VLIVE feeds from the database.",
		)
		return
	}

	if len(feeds) < 1 {
		ctx.RespondWarning(
			"This server has no VLIVE feeds set up.",
		)
		return
	}

	boardIDs := make([]int64, 0, len(feeds))
	for _, f := range feeds {
		boardIDs = append(boardIDs, f.BoardID)
	}
	slices.Compact(boardIDs)

	channelCodes := make([]string, len(boardIDs))
	for _, id := range boardIDs {
		board, err := db.VLIVE.GetBoard(id)
		if err != nil {
			continue
		}

		channelCodes = append(channelCodes, board.ChannelCode)
	}
	slices.Compact(channelCodes)

	boardTitles := make(map[int64]string)
	channelNames := make(map[string]string)
	for _, id := range boardIDs {
		boardTitles[id] = strconv.FormatInt(id, 10)
	}
	for _, code := range channelCodes {
		channel, res, err := vliveutil.GetChannel(code)
		if err != nil {
			log.Println(err)
		}
		if err != nil || res.StatusCode != http.StatusOK {
			channelNames[code] = code
			continue
		}

		channelNames[channel.Code] = channel.Name

		boards, res, err := vliveutil.GetUnwrappedBoards(channel.Code)
		if err != nil {
			log.Println(err)
		}
		if err != nil || res.StatusCode != http.StatusOK {
			continue
		}

		for _, b := range boards {
			boardTitles[b.ID] = b.Title
		}
	}

	feedList := make([]string, 0, len(feeds))
	for _, feed := range feeds {
		feedEntry := feed.ChannelID.Mention() + " | "

		channelName := "Unknown"
		boardTitle := "Unknown"
		board, err := db.VLIVE.GetBoard(feed.BoardID)
		if err != nil {
			log.Println(err)
		} else {
			name := channelNames[board.ChannelCode]
			if name != "" {
				channelName = name
			}
		}

		title := boardTitles[feed.BoardID]
		if title != "" {
			boardTitle = title
		}

		feedEntry += channelName + " - " + boardTitle

		feedFlags := make([]string, 0)
		if feed.PostTypes == vlivedb.VideosOnly {
			feedFlags = append(feedFlags, "📼")
		} else if feed.PostTypes == vlivedb.PostsOnly {
			feedFlags = append(feedFlags, "📝")
		}
		if feed.Reposts {
			feedFlags = append(feedFlags, "♻️")
		}

		if len(feedFlags) > 0 {
			feedEntry += " (" + strings.Join(feedFlags, " + ") + ")"
		}

		feedList = append(feedList, feedEntry)
	}

	descriptionPages := util.PagedLines(feedList, 2048, 20)
	messagePages := make([]router.MessagePage, len(descriptionPages))
	numOfFeeds := len(feeds)
	numOfFeedsFooter := fmt.Sprintf(
		"%d %s",
		numOfFeeds,
		util.Pluralise("Feed", numOfFeeds),
	)

	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		messagePages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Name: "VLIVE Notifications", Icon: vliveIcon,
					},
					Description: description,
					Color:       vliveColour,
					Footer: &discord.EmbedFooter{
						Text: dctools.SeparateEmbedFooter(
							pageID,
							numOfFeedsFooter,
						),
					},
				},
			},
		}
	}

	ctx.RespondPaging(messagePages)
}
