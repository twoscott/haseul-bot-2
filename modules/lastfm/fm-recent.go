package lastfm

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/shkh/lastfm-go/lastfm"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var fmRecentCommand = &router.Command{
	Name:      "recent",
	Aliases:   []string{"recents"},
	UseTyping: true,
	Run:       fmRecentRun,
}

func fmRecentRun(ctx router.CommandCtx, args []string) {
	limit := 100

	if len(args) > 0 && args[0] != "" {
		limit, _ = strconv.Atoi(args[0])
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

	res, ok := recentTracks(ctx, lfUser, limit)
	if !ok {
		return
	}
	if len(res.Tracks) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			fmt.Sprintf("%s has not listened to any music.", lfUser),
		)
		return
	}

	messagePages := recentEmbeds(ctx, *res)

	cmdutil.ReplyWithPaging(ctx, ctx.Msg, messagePages)
}

func recentEmbeds(
	ctx router.CommandCtx,
	recentTracks lastfm.UserGetRecentTracks) []router.MessagePage {

	tracks := recentTracks.Tracks
	images := tracks[0].Images
	lfUser := recentTracks.User
	totalScrobbles := humanize.Comma(int64(recentTracks.Total))

	nowPlaying, _ := strconv.ParseBool(tracks[0].NowPlaying)
	authorTitle := util.Possessive(lfUser) + " Recent Scrobbles"
	authorURL := getLibraryURL(lfUser)

	thumbnailURL := images[len(images)-1].Url
	if thumbnailURL == "" {
		thumbnailURL = getThumbURL(noAlbumHash)
	}
	footerText := fmt.Sprintf(
		"Total Scrobbles: %s%sPowered by Last.fm",
		totalScrobbles, dctools.EmbedFooterSep,
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
			line = fmt.Sprintf("\\▶ %s (Now)", lineTrack)
		} else {
			var timeAgoString string

			unixTime, err := strconv.ParseInt(track.Date.Uts, 10, 64)
			if err != nil {
				log.Println(err)
				timeAgoString = "N/A"
			} else {
				timestamp := time.Unix(unixTime, 0)
				timeAgoString = util.MaxTimeAgoString(timestamp)
			}

			line = fmt.Sprintf("%d. %s (%s)", i+1, lineTrack, timeAgoString)
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
						Text: pageID + dctools.EmbedFooterSep + footerText,
					},
				},
			},
		}
	}

	return messagePages
}
