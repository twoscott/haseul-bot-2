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

var lastFmTopArtistsCommand = &router.SubCommand{
	Name:        "artists",
	Description: "Displays your most scrobbled artists",
	Handler: &router.CommandHandler{
		Executor: lastFmTopArtistsExec,
		Defer:    true,
	},
	Options: []discord.CommandOptionValue{
		&discord.IntegerOption{
			OptionName:  "artists",
			Description: "The number of top artists to display for the user",
			Min:         option.NewInt(1),
			Max:         option.NewInt(1000),
		},
		&discord.IntegerOption{
			OptionName:  "period",
			Description: "The period of time to search for top artists within",
			Choices:     timePeriodChoices,
		},
	},
}

func lastFmTopArtistsExec(ctx router.CommandCtx) {
	artistCount, _ := ctx.Options.Find("artists").IntValue()
	if artistCount == 0 {
		artistCount = 10
	}

	periodOption, _ := ctx.Options.Find("period").IntValue()
	timeframe := lastFmPeriod(periodOption).Timeframe()

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

	res, err := getTopArtists(timeframe, lfUser, artistCount)
	if err != nil {
		errMsg := errorResponseMessage(err)
		ctx.RespondError(errMsg)
		return
	}

	if len(res.Artists) < 1 {
		ctx.RespondWarning(
			"You have not scrobbled any tracks on Last.fm.",
		)
		return
	}

	messagePages := topArtistsEmbeds(*res, *timeframe)

	ctx.RespondPaging(messagePages)
}

func getTopArtists(
	tf *timeframe, lfUser string, limit int64) (*lastfm.UserGetTopArtists, error) {

	res, err := lf.User.GetTopArtists(
		lastfm.P{"user": lfUser, "limit": limit, "period": tf.apiPeriod},
	)
	if err != nil {
		return nil, err
	}

	if res.User == "" {
		res.User = lfUser
	}

	if int64(len(res.Artists)) > limit {
		res.Artists = res.Artists[:limit]
	}

	return &res, nil
}

func topArtistsEmbeds(
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
		thumbnailURL = getImageURL(noArtistHash)
	} else {
		thumbnailURL = toImage(thumbnailURL)
	}

	footerText := dctools.SeparateEmbedFooter(
		fmt.Sprintf("Total Artists: %s", totalArtists),
		"Powered by Last.fm",
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
						Text: dctools.SeparateEmbedFooter(pageID, footerText),
					},
				},
			},
		}
	}

	return messagePages
}
