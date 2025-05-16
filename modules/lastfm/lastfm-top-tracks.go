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

var lastFMTopTracksCommand = &router.SubCommand{
	Name:        "tracks",
	Description: "Displays your most scrobbled tracks",
	Handler: &router.CommandHandler{
		Executor: lastFMTopTracksExec,
		Defer:    true,
	},
	Options: []discord.CommandOptionValue{
		&discord.IntegerOption{
			OptionName:  "tracks",
			Description: "The number of top tracks to display for the user",
			Min:         option.NewInt(1),
			Max:         option.NewInt(1000),
		},
		&discord.StringOption{
			OptionName:  "period",
			Description: "The period of time to search for top tracks within",
			Choices:     timePeriodChoices,
		},
	},
}

func lastFMTopTracksExec(ctx router.CommandCtx) {
	trackCount, _ := ctx.Options.Find("tracks").IntValue()
	if trackCount == 0 {
		trackCount = 10
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

	res, err := fm.User.TopTracks(lastfm.UserTopTracksParams{
		User:   fmUser,
		Limit:  uint(trackCount),
		Period: timeframe.apiPeriod,
	})
	if err != nil {
		log.Println(err)
		if msg, ok := errMessage(err); ok {
			ctx.RespondError(msg)
		} else {
			ctx.RespondError("Unable to fetch your top tracks from Last.fm.")
		}
		return
	}

	if len(res.Tracks) < 1 {
		m := "You have not scrobbled any tracks on Last.fm during '%s'."
		ctx.RespondWarning(fmt.Sprintf(m, timeframe.displayPeriod))
		return
	}

	if res.User != "" {
		fmUser = res.User
	}

	trackList := make([]string, 0, len(res.Tracks))
	for i, track := range res.Tracks {
		artistName := dctools.EscapeMarkdown(track.Artist.Name)
		trackName := dctools.EscapeMarkdown(track.Title)
		trackLink := dctools.Hyperlink(trackName, track.URL)
		scrobbles := util.PluraliseWithCount("Scrobble", int64(track.Playcount))

		line := fmt.Sprintf("%d. %s - %s (%s)", i+1, artistName, trackLink, scrobbles)
		trackList = append(trackList, line)
	}

	title := util.Possessive(fmUser) + " Top Tracks"

	firstArtistName := res.Tracks[0].Artist.Name
	imageURL, err := scrapeArtistImage(firstArtistName, lastfm.ImgSizeLarge)
	if err != nil {
		log.Println(err)
		imageURL = lastfm.NoTrackImageURL.Resize(lastfm.ImgSizeLarge)
	}

	footer := util.PluraliseWithCount("Track", int64(res.Total))
	footer += " " + timeframe.displayPeriod
	footer = dctools.SeparateEmbedFooter(footer, "Powered by Last.fm")

	pages := util.PagedLines(trackList, 2048, 25)
	messagePages := make([]router.MessagePage, len(pages))
	for i, page := range pages {
		id := fmt.Sprintf("Page %d/%d", i+1, len(pages))

		e := discord.Embed{
			Title: timeframe.displayPeriod,
			Author: &discord.EmbedAuthor{
				Name: title,
				URL:  trackLibraryURL(fmUser, *timeframe),
				Icon: trackIcon,
			},
			Description: page,
			Color:       trackColour,
			Thumbnail:   &discord.EmbedThumbnail{URL: imageURL},
			Footer: &discord.EmbedFooter{
				Text: dctools.SeparateEmbedFooter(id, footer),
			},
		}

		messagePages[i] = router.MessagePage{Embeds: []discord.Embed{e}}
	}

	ctx.RespondPaging(messagePages)
}
