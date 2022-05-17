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

var fmTaCommand = &router.Command{
	Name:      "ta",
	Aliases:   []string{"topartists"},
	UseTyping: true,
	Run:       fmTaRun,
}

func fmTaRun(ctx router.CommandCtx, args []string) {
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

	res, ok := getTopArtists(ctx, timeframe, lfUser, limit)
	if !ok {
		return
	}
	if len(res.Artists) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			fmt.Sprintf("%s has not listened to any music.", lfUser),
		)
		return
	}

	messagePages := topArtistsEmbeds(ctx, *res, *timeframe)

	cmdutil.ReplyWithPaging(ctx, ctx.Msg, messagePages)
}

func getTopArtists(
	ctx router.CommandCtx,
	tf *timeframe,
	lfUser string,
	limit int) (*lastfm.UserGetTopArtists, bool) {

	res, err := lf.User.GetTopArtists(
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

	if len(res.Artists) > limit {
		res.Artists = res.Artists[:limit]
	}

	return &res, true
}

func topArtistsEmbeds(
	ctx router.CommandCtx,
	topArtists lastfm.UserGetTopArtists,
	tf timeframe) []router.MessagePage {

	artists := topArtists.Artists
	lfUser := topArtists.User
	totalArtists := humanize.Comma(int64(topArtists.Total))

	authorTitle := util.Possessive(lfUser) + " Top Artists"
	authorURL := getArtistLibraryURL(lfUser, tf)
	title := tf.displayPeriod

	thumbnailURL, err := scrapeArtistImage(artists[0].Name)
	if err != nil {
		thumbnailURL = getThumbURL(noArtistHash)
	}
	footerText := fmt.Sprintf(
		"Total Artists: %s%sPowered by Last.fm",
		totalArtists, dctools.EmbedFooterSep,
	)

	artistList := make([]string, 0, len(artists))
	for i, artist := range artists {
		var playCount string

		int64Plays, err := strconv.ParseInt(artist.PlayCount, 10, 64)
		if err != nil {
			playCount = "N/A"
		} else {
			playCount = humanize.Comma(int64Plays)
		}

		artistName := dctools.EscapeMarkdown(artist.Name)
		line := fmt.Sprintf(
			"%d. %s (%s Scrobbles)",
			i+1, dctools.Hyperlink(artistName, artist.Url), playCount,
		)

		artistList = append(artistList, line)
	}

	descriptionPages := util.PagedLines(artistList, 2048, 25)
	messagePages := make([]router.MessagePage, len(descriptionPages))
	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		messagePages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Name: authorTitle, URL: authorURL, Icon: artistIcon,
					},
					Title:       title,
					Description: description,
					Thumbnail:   &discord.EmbedThumbnail{URL: thumbnailURL},
					Color:       artistColour,
					Footer: &discord.EmbedFooter{
						Text: pageID + dctools.EmbedFooterSep + footerText,
					},
				},
			},
		}
	}

	return messagePages
}
