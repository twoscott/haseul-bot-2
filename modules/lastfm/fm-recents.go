package lastfm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/dustin/go-humanize"
	"github.com/shkh/lastfm-go/lastfm"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var fmRecentCommand = &router.SubCommand{
	Name:        "recents",
	Description: "Displays your recently scrobbled tracks",
	Handler: &router.CommandHandler{
		Executor: fmRecentExec,
		Defer:    true,
	},
	Options: []discord.CommandOptionValue{
		&discord.IntegerOption{
			OptionName:  "tracks",
			Description: "The number of recent tracks to display for the user",
			MinValue:    option.NewInt(1),
			MaxValue:    option.NewInt(1000),
		},
	},
}

func fmRecentExec(ctx router.CommandCtx) {
	trackCount, _ := ctx.Options.Find("tracks").IntValue()
	if trackCount < 1 {
		trackCount = 10
	}

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

	res, err := getRecentTracks(lfUser, trackCount)
	if err != nil {
		errMsg := errorResponseMessage(err)
		ctx.RespondError(errMsg)
		return
	}

	messagePages := recentsListEmbeds(res)

	ctx.RespondPaging(messagePages)
}

func recentsListEmbeds(
	recentTracks *lastfm.UserGetRecentTracks) []router.MessagePage {

	tracks := recentTracks.Tracks
	images := tracks[0].Images
	lfUser := recentTracks.User
	totalScrobbles := humanize.Comma(int64(recentTracks.Total))

	nowPlaying, _ := strconv.ParseBool(tracks[0].NowPlaying)
	authorTitle := util.Possessive(lfUser) + " Recent Scrobbles"
	authorURL := getLibraryURL(lfUser)

	thumbnailURL := images[len(images)-1].Url
	if thumbnailURL == "" {
		thumbnailURL = getImageURL(noAlbumHash)
	} else {
		thumbnailURL = toImage(thumbnailURL)
	}

	footerText := dctools.SeparateEmbedFooter(
		fmt.Sprintf("Total Scrobbles: %s", totalScrobbles),
		"Powered by Last.fm",
	)

	trackList := make([]string, 0, len(recentTracks.Tracks))
	for i, track := range recentTracks.Tracks {
		trackElems := dctools.MultiEscapeMarkdown(track.Artist.Name, track.Name)
		lineTrack := fmt.Sprintf(
			"%s - %s",
			trackElems[0], dctools.Hyperlink(trackElems[1], track.Url),
		)

		var line string
		if i == 0 && nowPlaying {
			line = fmt.Sprintf("1. %s", lineTrack)
		} else {
			var timeAgoString string

			unixTime, err := strconv.ParseInt(track.Date.Uts, 10, 64)
			if err != nil {
				log.Println(err)
				timeAgoString = "N/A"
			} else {
				timestamp := time.Unix(unixTime, 0)
				timeAgoString = dctools.UnixTimestampStyled(
					timestamp, dctools.RelativeTime,
				)
			}

			line = fmt.Sprintf("%d. %s %s", i+1, lineTrack, timeAgoString)
		}

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
						Name: authorTitle, URL: authorURL, Icon: scrobbleIcon,
					},
					Description: description,
					Thumbnail:   &discord.EmbedThumbnail{URL: thumbnailURL},
					Color:       scrobbleColour,
					Footer: &discord.EmbedFooter{
						Text: dctools.SeparateEmbedFooter(pageID, footerText),
					},
				},
			},
		}
	}

	return messagePages
}
