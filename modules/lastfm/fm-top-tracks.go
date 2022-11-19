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

var fmTopTracksCommand = &router.SubCommand{
	Name:        "tracks",
	Description: "Displays your most scrobbled tracks",
	Handler: &router.CommandHandler{
		Executor: fmTopTracksExec,
		Defer:    true,
	},
	Options: []discord.CommandOptionValue{
		&discord.IntegerOption{
			OptionName:  "tracks",
			Description: "The number of top tracks to display for the user",
			MinValue:    option.NewInt(1),
			MaxValue:    option.NewInt(1000),
		},
		&discord.IntegerOption{
			OptionName:  "period",
			Description: "The period of time to search for top tracks within",
			Choices:     timePeriodChoices,
		},
	},
}

func fmTopTracksExec(ctx router.CommandCtx) {
	trackCount, _ := ctx.Options.Find("tracks").IntValue()
	if trackCount == 0 {
		trackCount = 10
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

	res, err := getTopTracks(timeframe, lfUser, trackCount)
	if err != nil {
		errMsg := errorResponseMessage(err)
		ctx.RespondError(errMsg)
		return
	}

	if len(res.Tracks) < 1 {
		ctx.RespondWarning(
			"You have not scrobbled any tracks on Last.fm",
		)
		return
	}

	messagePages := topTracksEmbeds(*res, *timeframe)

	err = ctx.RespondPaging(messagePages)
	if err != nil {
		log.Println(err)
	}
}

func getTopTracks(
	tf *timeframe,
	lfUser string,
	limit int64) (*lastfm.UserGetTopTracks, error) {

	res, err := lf.User.GetTopTracks(
		lastfm.P{"user": lfUser, "limit": limit, "period": tf.apiPeriod},
	)
	if err != nil {
		return nil, err
	}

	if res.User == "" {
		res.User = lfUser
	}

	if int64(len(res.Tracks)) > limit {
		res.Tracks = res.Tracks[:limit]
	}

	return &res, nil
}

func topTracksEmbeds(
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
		thumbnailURL = getImageURL(noTrackHash)
	} else {
		thumbnailURL = toImage(thumbnailURL)
	}

	footerText := dctools.SeparateEmbedFooter(
		fmt.Sprintf("Total Tracks: %s", totalTracks),
		"Powered by Last.fm",
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
						Text: dctools.SeparateEmbedFooter(pageID, footerText),
					},
				},
			},
		}
	}

	return messagePages
}
