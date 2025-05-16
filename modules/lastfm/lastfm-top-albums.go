package lastfm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"

	"github.com/twoscott/gobble-fm/lastfm"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var lastFMTopAlbumsCommand = &router.SubCommand{
	Name:        "albums",
	Description: "Displays your most scrobbled albums",
	Handler: &router.CommandHandler{
		Executor: lastFMTopAlbumsExec,
		Defer:    true,
	},
	Options: []discord.CommandOptionValue{
		&discord.IntegerOption{
			OptionName:  "albums",
			Description: "The number of top albums to display for the user",
			Min:         option.NewInt(1),
			Max:         option.NewInt(1000),
		},
		&discord.StringOption{
			OptionName:  "period",
			Description: "The period of time to search for top albums within",
			Choices:     timePeriodChoices,
		},
	},
}

func lastFMTopAlbumsExec(ctx router.CommandCtx) {
	albumCount, _ := ctx.Options.Find("albums").IntValue()
	if albumCount == 0 {
		albumCount = 10
	}

	periodOption := ctx.Options.Find("period").String()
	timeframe := newTimeframe(lastfm.Period(periodOption))

	fmUser, err := db.LastFM.GetUser(ctx.Interaction.SenderID())
	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondWarning(
			fmt.Sprintf(
				"Please link a Last.fm username to your account using %s",
				lastFMSetCommand.Mention(),
			),
		)
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondGenericError()
		return
	}

	res, err := fm.User.TopAlbums(lastfm.UserTopAlbumsParams{
		User:   fmUser,
		Limit:  uint(albumCount),
		Period: timeframe.apiPeriod,
	})
	if err != nil {
		log.Println(err)
		if msg, ok := errMessage(err); ok {
			ctx.RespondError(msg)
		} else {
			ctx.RespondError("Unable to fetch your top albums from Last.fm.")
		}
		return
	}

	if len(res.Albums) < 1 {
		m := "You have not scrobbled any albums on Last.fm during '%s'."
		ctx.RespondWarning(fmt.Sprintf(m, timeframe.displayPeriod))
		return
	}

	if res.User != "" {
		fmUser = res.User
	}

	albumList := make([]string, 0, len(res.Albums))
	for i, album := range res.Albums {
		artistName := dctools.EscapeMarkdown(album.Artist.Name)
		albumName := dctools.EscapeMarkdown(album.Title)
		albumLink := dctools.Hyperlink(albumName, album.URL)
		scrobbles := util.PluraliseWithCount("Scrobble", int64(album.Playcount))

		line := fmt.Sprintf("%d. %s - %s (%s)", i+1, artistName, albumLink, scrobbles)
		albumList = append(albumList, line)
	}

	title := util.Possessive(fmUser) + " Top Albums"

	imageURL := res.Albums[0].Cover.SizedURL(lastfm.ImgSizeLarge)
	if imageURL == "" {
		imageURL = lastfm.NoAlbumImageURL.Resize(lastfm.ImgSizeLarge)
	}

	footer := util.PluraliseWithCount("Album", int64(res.Total))
	footer += " " + timeframe.displayPeriod
	footer = dctools.SeparateEmbedFooter(footer, "Powered by Last.fm")

	pages := util.PagedLines(albumList, 2048, 25)
	messagePages := make([]router.MessagePage, len(pages))
	for i, page := range pages {
		id := fmt.Sprintf("Page %d/%d", i+1, len(pages))

		e := discord.Embed{
			Title: timeframe.displayPeriod,
			Author: &discord.EmbedAuthor{
				Name: title,
				URL:  albumLibraryURL(fmUser, *timeframe),
				Icon: albumIcon,
			},
			Description: page,
			Color:       albumColour,
			Thumbnail:   &discord.EmbedThumbnail{URL: imageURL},
			Footer: &discord.EmbedFooter{
				Text: dctools.SeparateEmbedFooter(id, footer),
			},
		}

		messagePages[i] = router.MessagePage{Embeds: []discord.Embed{e}}
	}

	ctx.RespondPaging(messagePages)
}
