package lastfm

import (
	"database/sql"
	"errors"
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

var lastFmCurrentCommand = &router.SubCommand{
	Name:        "current",
	Description: "Displays your currently scrobbling track",
	Handler: &router.CommandHandler{
		Executor: lastFmCurrentExec,
	},
}

func lastFmCurrentExec(ctx router.CommandCtx) {
	lfUser, err := db.LastFM.GetUser(ctx.Interaction.SenderID())

	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondWarning(
			fmt.Sprintf(
				"Please link a Last.fm username to your account using %s",
				lastFmSetCommand.Mention(),
			),
		)
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondGenericError()
		return
	}

	res, err := getRecentTracks(lfUser, 1)
	if err != nil {
		errMsg := errorResponseMessage(err)
		ctx.RespondError(errMsg)
		return
	}

	if len(res.Tracks) < 1 {
		ctx.RespondWarning(
			"You have not scrobbled any tracks on Last.fm.",
		)
		return
	}

	var playCount string
	var userLoved string
	track, err := getTrackInfo(res)
	if err != nil {
		playCount = "N/A"
		userLoved = "0"
	} else {
		playCount = track.UserPlayCount
		userLoved = track.UserLoved
	}

	embed := npEmbed(ctx, *res, playCount, userLoved)

	ctx.RespondEmbed(*embed)
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
		thumbnailURL = getImageURL(noAlbumHash)
	} else {
		thumbnailURL = toImage(thumbnailURL)
	}

	trackFieldName := "Song"
	trackFieldElems := dctools.MultiEscapeMarkdown(track.Artist.Name, track.Name)
	trackFieldValue := fmt.Sprintf(
		"%s - %s",
		trackFieldElems[0], dctools.Hyperlink(trackFieldElems[1], track.Url),
	)

	if playCount == "0" {
		playCount = "First"
	}
	footerText := dctools.SeparateEmbedFooter(
		fmt.Sprintf("Song Scrobbles: %s", playCount),
		"Powered by Last.fm",
	)

	if loved {
		lovedText := "❤ Loved"
		footerText = dctools.SeparateEmbedFooter(footerText, lovedText)
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
