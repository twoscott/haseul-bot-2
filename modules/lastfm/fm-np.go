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

var fmNpCommand = &router.Command{
	Name:      "np",
	Aliases:   []string{"nowplaying"},
	UseTyping: true,
	Run:       fmNpRun,
}

func fmNpRun(ctx router.CommandCtx, _ []string) {
	lfUser, ok := getLfUser(ctx)
	if !ok {
		return
	}

	res, ok := recentTracks(ctx, lfUser, 1)
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
	var userLoved string
	track, ok := trackInfo(ctx, res)
	if !ok {
		playCount = "N/A"
		userLoved = "0"
	} else {
		playCount = track.UserPlayCount
		userLoved = track.UserLoved
	}

	embed := npEmbed(ctx, *res, playCount, userLoved)

	dctools.EmbedReplyNoPing(ctx.State, ctx.Msg, *embed)
}

func npEmbed(
	ctx router.CommandCtx,
	recentTracks lastfm.UserGetRecentTracks,
	playCount string,
	userLoved string) *discord.Embed {

	track := &recentTracks.Tracks[0]
	images := track.Images
	lfUser := recentTracks.User
	loved, _ := strconv.ParseBool(userLoved)

	nowPlaying, _ := strconv.ParseBool(track.NowPlaying)
	authorTitle := util.Possessive(lfUser)
	if nowPlaying {
		authorTitle += " Now Scrobbling"
	} else {
		authorTitle += " Last Scrobbled"
	}
	authorURL := lastFmURL + "/user/" + lfUser + "/library"

	thumbnailURL := images[len(images)-1].Url
	if thumbnailURL == "" {
		thumbnailURL = getThumbURL(noAlbumHash)
	}

	trackFieldName := "Song"
	trackFieldElems := dctools.MultiEscapeMarkdown(track.Artist.Name, track.Name)
	trackFieldValue := fmt.Sprintf(
		"%s - %s",
		trackFieldElems[0], dctools.Hyperlink(trackFieldElems[1], track.Url),
	)

	footerText := fmt.Sprintf(
		"Song Scrobbles: %s%sPowered by Last.fm",
		playCount, dctools.EmbedFooterSep,
	)
	if loved {
		lovedText := "❤ Loved"
		footerText = lovedText + dctools.EmbedFooterSep + footerText
	}

	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: authorTitle, URL: authorURL, Icon: scrobbleIcon,
		},
		Thumbnail: &discord.EmbedThumbnail{URL: thumbnailURL},
		Fields: []discord.EmbedField{
			{Name: trackFieldName, Value: trackFieldValue, Inline: false},
		},
		Color:  scrobbleColour,
		Footer: &discord.EmbedFooter{Text: footerText},
	}

	if track.Album.Name != "" {
		albumFieldName := "Album"
		albumFieldValue := dctools.EscapeMarkdown(track.Album.Name)
		albumField := discord.EmbedField{
			Name: albumFieldName, Value: albumFieldValue, Inline: false,
		}
		embed.Fields = append(embed.Fields, albumField)
	}

	if !nowPlaying {
		timestamp, err := strconv.ParseInt(track.Date.Uts, 10, 64)
		if err != nil {
			log.Println(err)
		} else {
			time := time.Unix(timestamp, 0)
			embed.Timestamp = discord.NewTimestamp(time)
		}
	}

	return &embed
}
