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

var fmTalbCommand = &router.Command{
	Name:      "talb",
	Aliases:   []string{"topalbums"},
	UseTyping: true,
	Run:       fmTalbRun,
}

func fmTalbRun(ctx router.CommandCtx, args []string) {
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

	res, ok := getTopAlbums(ctx, timeframe, lfUser, limit)
	if !ok {
		return
	}
	if len(res.Albums) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			fmt.Sprintf("%s has not listened to any music.", lfUser),
		)
		return
	}

	messagePages := topAlbumsEmbeds(ctx, *res, *timeframe)

	cmdutil.ReplyWithPaging(ctx, ctx.Msg, messagePages)
}

func getTopAlbums(
	ctx router.CommandCtx,
	tf *timeframe,
	lfUser string,
	limit int) (*lastfm.UserGetTopAlbums, bool) {

	res, err := lf.User.GetTopAlbums(
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

	if len(res.Albums) > limit {
		res.Albums = res.Albums[:limit]
	}

	return &res, true
}

func topAlbumsEmbeds(
	ctx router.CommandCtx,
	topAlbums lastfm.UserGetTopAlbums,
	tf timeframe) []router.MessagePage {

	albums := topAlbums.Albums
	images := albums[0].Images
	lfUser := topAlbums.User
	totalAlbums := humanize.Comma(int64(topAlbums.Total))

	authorTitle := util.Possessive(lfUser) + " Top Albums"
	authorURL := getAlbumLibraryURL(lfUser, tf)
	title := tf.displayPeriod

	thumbnailURL := images[len(images)-1].Url
	if thumbnailURL == "" {
		thumbnailURL = getThumbURL(noAlbumHash)
	}
	footerText := fmt.Sprintf(
		"Total Albums: %s%sPowered by Last.fm",
		totalAlbums, dctools.EmbedFooterSep,
	)

	albumList := make([]string, 0, len(albums))
	for i, album := range albums {
		var playCount string

		int64Plays, err := strconv.ParseInt(album.PlayCount, 10, 64)
		if err != nil {
			playCount = "N/A"
		} else {
			playCount = humanize.Comma(int64Plays)
		}

		albumElems := dctools.MultiEscapeMarkdown(album.Artist.Name, album.Name)
		line := fmt.Sprintf(
			"%d. %s - %s (%s Scrobbles)",
			i+1, albumElems[0],
			dctools.Hyperlink(albumElems[1], album.Url),
			playCount,
		)

		albumList = append(albumList, line)
	}

	descriptionPages := util.PagedLines(albumList, 2048, 25)
	messagePages := make([]router.MessagePage, len(descriptionPages))
	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		messagePages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Name: authorTitle, URL: authorURL, Icon: albumIcon,
					},
					Title:       title,
					Description: description,
					Thumbnail:   &discord.EmbedThumbnail{URL: thumbnailURL},
					Color:       albumColour,
					Footer: &discord.EmbedFooter{
						Text: pageID + dctools.EmbedFooterSep + footerText,
					},
				},
			},
		}
	}

	return messagePages
}
