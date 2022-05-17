package lastfm

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/shkh/lastfm-go/lastfm"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var fmCommand = &router.Command{
	Name:        "fm",
	Aliases:     []string{"lf", "lastfm"},
	UseTyping:   true,
	Run:         fmRun,
	SubCommands: make(router.CommandMap),
}

func fmRun(ctx router.CommandCtx, _ []string) {
	lfUser, ok := getLfUser(ctx)
	if !ok {
		return
	}

	res, ok := recentTracks(ctx, lfUser, 2)
	if !ok {
		return
	}
	if len(res.Tracks) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			fmt.Sprintf("%s has not listened to any music.", lfUser),
		)
		return
	}

	var playCount string
	track, ok := trackInfo(ctx, res)
	if !ok {
		playCount = "N/A"
	} else {
		playCount = track.UserPlayCount
	}

	fmEmbed(ctx, res, playCount)
}

func fmEmbed(
	ctx router.CommandCtx,
	recentTracks *lastfm.UserGetRecentTracks,
	playCount string) {

	track1 := &recentTracks.Tracks[0]
	track2 := &recentTracks.Tracks[1]
	images := track1.Images
	lfUser := recentTracks.User

	nowPlaying, _ := strconv.ParseBool(track1.NowPlaying)
	authorTitle := util.Possessive(recentTracks.User) + " Recent Scrobbles"
	authorURL := lastFmURL + "/user/" + lfUser + "/library"

	thumbnailURL := images[len(images)-1].Url
	if thumbnailURL == "" {
		thumbnailURL = getThumbURL(noAlbumHash)
	}

	var field1Name string
	if nowPlaying {
		field1Name = "Now Scrobbling"
	} else {
		field1Name = "Last Scrobbled"
	}
	field2Name := "Previously Scrobbled"

	field1Elems := dctools.MultiEscapeMarkdown(track1.Artist.Name, track1.Name)
	field2Elems := dctools.MultiEscapeMarkdown(track2.Artist.Name, track2.Name)
	field1Value := fmt.Sprintf(
		"%s - %s",
		field1Elems[0], dctools.Hyperlink(field1Elems[1], track1.Url),
	)
	field2Value := fmt.Sprintf(
		"%s - %s",
		field2Elems[0], dctools.Hyperlink(field2Elems[1], track1.Url),
	)

	if track1.Album.Name != "" {
		field1Value += fmt.Sprintf(
			" | %s", dctools.EscapeMarkdown(track1.Album.Name),
		)
	}
	if track2.Album.Name != "" {
		field2Value += fmt.Sprintf(
			" | %s", dctools.EscapeMarkdown(track2.Album.Name),
		)
	}

	footerText := fmt.Sprintf(
		"Song Scrobbles: %s%sPowered by Last.fm",
		playCount, dctools.EmbedFooterSep,
	)

	embed := &discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: authorTitle, URL: authorURL, Icon: scrobbleIcon,
		},
		Thumbnail: &discord.EmbedThumbnail{URL: thumbnailURL},
		Fields: []discord.EmbedField{
			{Name: field1Name, Value: field1Value, Inline: false},
			{Name: field2Name, Value: field2Value, Inline: false},
		},
		Color:  scrobbleColour,
		Footer: &discord.EmbedFooter{Text: footerText},
	}

	if !nowPlaying {
		timestamp, err := strconv.ParseInt(track1.Date.Uts, 10, 64)
		if err != nil {
			log.Println(err)
		} else {
			time := time.Unix(timestamp, 0)
			embed.Timestamp = discord.NewTimestamp(time)
		}
	}

	dctools.EmbedReplyNoPing(ctx.State, ctx.Msg, *embed)
}
