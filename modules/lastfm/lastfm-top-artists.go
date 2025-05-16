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

var lastFMTopArtistsCommand = &router.SubCommand{
	Name:        "artists",
	Description: "Displays your most scrobbled artists",
	Handler: &router.CommandHandler{
		Executor: lastFMTopArtistsExec,
		Defer:    true,
	},
	Options: []discord.CommandOptionValue{
		&discord.IntegerOption{
			OptionName:  "artists",
			Description: "The number of top artists to display for the user",
			Min:         option.NewInt(1),
			Max:         option.NewInt(1000),
		},
		&discord.StringOption{
			OptionName:  "period",
			Description: "The period of time to search for top artists within",
			Choices:     timePeriodChoices,
		},
	},
}

func lastFMTopArtistsExec(ctx router.CommandCtx) {
	artistCount, _ := ctx.Options.Find("artists").IntValue()
	if artistCount == 0 {
		artistCount = 10
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

	res, err := fm.User.TopArtists(lastfm.UserTopArtistsParams{
		User:   fmUser,
		Limit:  uint(artistCount),
		Period: timeframe.apiPeriod,
	})
	if err != nil {
		log.Println(err)
		if msg, ok := errMessage(err); ok {
			ctx.RespondError(msg)
		} else {
			ctx.RespondError("Unable to fetch your top artists from Last.fm.")
		}
		return
	}

	if len(res.Artists) < 1 {
		m := "You have not scrobbled any artists on Last.fm during '%s'."
		ctx.RespondWarning(fmt.Sprintf(m, timeframe.displayPeriod))
		return
	}

	if res.User != "" {
		fmUser = res.User
	}

	artistList := make([]string, 0, len(res.Artists))
	for i, artist := range res.Artists {
		artistName := dctools.EscapeMarkdown(artist.Name)
		artistLink := dctools.Hyperlink(artistName, artist.URL)
		scrobbles := util.PluraliseWithCount("Scrobble", int64(artist.Playcount))

		line := fmt.Sprintf("%d. %s (%s)", i+1, artistLink, scrobbles)
		artistList = append(artistList, line)
	}

	title := util.Possessive(fmUser) + " Top Artists"

	firstName := res.Artists[0].Name
	imageURL, err := scrapeArtistImage(firstName, lastfm.ImgSizeLarge)
	if err != nil {
		log.Println(err)
		imageURL = lastfm.NoArtistImageURL.Resize(lastfm.ImgSizeLarge)
	}

	footer := util.PluraliseWithCount("Artist", int64(res.Total))
	footer += " " + timeframe.displayPeriod
	footer = dctools.SeparateEmbedFooter(footer, "Powered by Last.fm")

	pages := util.PagedLines(artistList, 2048, 25)
	messagePages := make([]router.MessagePage, len(pages))
	for i, page := range pages {
		id := fmt.Sprintf("Page %d/%d", i+1, len(pages))

		e := discord.Embed{
			Title: timeframe.displayPeriod,
			Author: &discord.EmbedAuthor{
				Name: title,
				URL:  artistLibraryURL(fmUser, *timeframe),
				Icon: artistIcon,
			},
			Description: page,
			Color:       artistColour,
			Thumbnail:   &discord.EmbedThumbnail{URL: imageURL},
			Footer: &discord.EmbedFooter{
				Text: dctools.SeparateEmbedFooter(id, footer),
			},
		}

		messagePages[i] = router.MessagePage{Embeds: []discord.Embed{e}}
	}

	ctx.RespondPaging(messagePages)
}
