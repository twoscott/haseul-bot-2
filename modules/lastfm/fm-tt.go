package lastfm

import (
	"fmt"
	"log"
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/shkh/lastfm-go/lastfm"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var fmTtCommand = &router.Command{
	Name:      "tt",
	Aliases:   []string{"toptracks"},
	UseTyping: true,
	Run:       fmTtRun,
}

func fmTtRun(ctx router.CommandCtx, args []string) {
	timeframe := getTimeframe("7day")
	limit := 100

	if len(args) > 0 && args[0] != "" {
		timeframe = getTimeframe(args[0])
	}
	if len(args) > 1 && args[1] != "" {
		limit, _ = strconv.Atoi(args[1])
	}

	if limit < 1 {
		limit = 1
	} else if limit > 1000 {
		limit = 1000
	}

	lfUser, ok := getLfUser(ctx)
	if !ok {
		return
	}

	res, ok := getTopTracks(ctx, timeframe, lfUser, limit)
	if !ok {
		return
	}
	if len(res.Tracks) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			fmt.Sprintf("%s has not listened to any music.", lfUser),
		)
		return
	}

	messagePages := topTracksEmbeds(ctx, *res, *timeframe)

	cmdutil.ReplyWithPaging(ctx, ctx.Msg, messagePages)
}

func getTopTracks(ctx router.CommandCtx,
	tf *timeframe,
	lfUser string,
	limit int) (*lastfm.UserGetTopTracks, bool) {

	res, err := lf.User.GetTopTracks(
		lastfm.P{"user": lfUser, "limit": limit, "period": tf.apiPeriod},
	)
	if err != nil {
		lfErr := getLfError(err)
		switch lfErr.Code {
		case 6:
			dctools.ReplyWithWarning(ctx.State, ctx.Msg,
				fmt.Sprintf("Last.fm user %s does not exist.", lfUser),
			)
		case 8:
			log.Println(err)
			dctools.ReplyWithError(ctx.State, ctx.Msg,
				"I could not get a response from Last.fm. Please try again.",
			)
		default:
			log.Println(err)
			dctools.ReplyWithError(ctx.State, ctx.Msg,
				"Unknown Last.fm Error occurred.",
			)
		}
		return nil, false
	}

	if res.User == "" {
		res.User = lfUser
	}

	if len(res.Tracks) > limit {
		res.Tracks = res.Tracks[:limit]
	}

	return &res, true
}

func topTracksEmbeds(
	ctx router.CommandCtx,
	topTracks lastfm.UserGetTopTracks,
	tf timeframe) []router.MessagePage {

	tracks := topTracks.Tracks
	lfUser := topTracks.User
	totalTracks := humanize.Comma(int64(topTracks.Total))

	authorTitle := util.Possessive(lfUser) + " Top Tracks"
	authorURL := getTrackLibraryURL(lfUser, tf)
	title := tf.displayPeriod

	thumbnailURL, err := scrapeArtistImage(tracks[0].Artist.Name)
	if err != nil {
		thumbnailURL = getThumbURL(noArtistHash)
	}
	footerText := fmt.Sprintf(
		"Total Tracks: %s%sPowered by Last.fm",
		totalTracks, dctools.EmbedFooterSep,
	)

	trackList := make([]string, 0, len(tracks))
	for i, track := range tracks {
		var playCount string

		int64Plays, err := strconv.ParseInt(track.PlayCount, 10, 64)
		if err != nil {
			playCount = "N/A"
		} else {
			playCount = humanize.Comma(int64Plays)
		}

		trackElems := dctools.MultiEscapeMarkdown(track.Artist.Name, track.Name)
		line := fmt.Sprintf(
			"%d. %s - %s (%s Scrobbles)",
			i+1, trackElems[0],
			dctools.Hyperlink(trackElems[1], track.Url),
			playCount,
		)

		trackList = append(trackList, line)
	}

	descriptionPages := util.PagedLines(trackList, 2048, 25)
	messagePages := make([]router.MessagePage, len(descriptionPages))
	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		messagePages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Name: authorTitle, URL: authorURL, Icon: trackIcon,
					},
					Title:       title,
					Description: description,
					Thumbnail:   &discord.EmbedThumbnail{URL: thumbnailURL},
					Color:       trackColour,
					Footer: &discord.EmbedFooter{
						Text: pageID + dctools.EmbedFooterSep + footerText,
					},
				},
			},
		}
	}

	return messagePages
}
