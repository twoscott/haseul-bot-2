package lastfm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/dustin/go-humanize"
	"github.com/shkh/lastfm-go/lastfm"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var fmTopAlbumsCommand = &router.SubCommand{
	Name:        "albums",
	Description: "Displays your most scrobbled albums",
	Handler: &router.CommandHandler{
		Executor: fmTopAlbumsExec,
		Defer:    true,
	},
	Options: []discord.CommandOptionValue{
		&discord.IntegerOption{
			OptionName:  "albums",
			Description: "The number of top albums to display for the user",
			MinValue:    option.NewInt(1),
			MaxValue:    option.NewInt(1000),
		},
		&discord.IntegerOption{
			OptionName:  "period",
			Description: "The period of time to search for top albums within",
			Choices:     timePeriodChoices,
		},
	},
}

func fmTopAlbumsExec(ctx router.CommandCtx) {
	albumCount, _ := ctx.Options.Find("albums").IntValue()
	if albumCount == 0 {
		albumCount = 10
	}

	periodOption, _ := ctx.Options.Find("period").IntValue()
	period := lastFmPeriod(periodOption)
	timeframe := period.Timeframe()

	lfUser, err := db.LastFM.GetUser(ctx.Interaction.SenderID())
	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondWarning(
			"Please link a Last.fm username to your account using `/fm set`",
		)
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondGenericError()
		return
	}

	res, err := getTopAlbums(timeframe, lfUser, albumCount)
	if err != nil {
		errMsg := errorResponseMessage(err)
		ctx.RespondError(errMsg)
		return
	}

	if len(res.Albums) < 1 {
		ctx.RespondWarning(
			"You have not scrobbled any tracks on Last.fm",
		)
		return
	}

	messagePages := topAlbumsEmbeds(*res, *timeframe)

	ctx.RespondPaging(messagePages)
}

func getTopAlbums(
	tf *timeframe,
	lfUser string,
	limit int64) (*lastfm.UserGetTopAlbums, error) {

	res, err := lf.User.GetTopAlbums(
		lastfm.P{"user": lfUser, "limit": limit, "period": tf.apiPeriod},
	)
	if err != nil {
		return nil, err
	}

	if res.User == "" {
		res.User = lfUser
	}

	if int64(len(res.Albums)) > limit {
		res.Albums = res.Albums[:limit]
	}

	return &res, nil
}

func topAlbumsEmbeds(
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
		thumbnailURL = getImageURL(noAlbumHash)
	} else {
		thumbnailURL = toImage(thumbnailURL)
	}

	footerText := dctools.SeparateEmbedFooter(
		fmt.Sprintf("Total Albums: %s", totalAlbums),
		"Powered by Last.fm",
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
						Text: dctools.SeparateEmbedFooter(pageID, footerText),
					},
				},
			},
		}
	}

	return messagePages
}
