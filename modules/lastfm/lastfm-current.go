package lastfm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/gobble-fm/lastfm"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var lastFMCurrentCommand = &router.SubCommand{
	Name:        "current",
	Description: "Displays your currently scrobbling track",
	Handler: &router.CommandHandler{
		Executor: lastFMCurrentExec,
		Defer:    true,
	},
}

func lastFMCurrentExec(ctx router.CommandCtx) {
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

	res, err := fm.User.RecentTrack(fmUser)
	if err != nil {
		log.Println(err)
		if msg, ok := errMessage(err); ok {
			ctx.RespondError(msg)
		} else {
			ctx.RespondError("Unable to fetch your recent scrobble from Last.fm.")
		}
		return
	}

	if res.Track == nil {
		ctx.RespondWarning("You have not scrobbled any tracks on Last.fm.")
		return
	}

	track := res.Track

	var (
		playcount int
		loved     bool
	)

	info, err := fm.Track.UserInfo(lastfm.TrackUserInfoParams{
		Artist: track.Artist.Name,
		Track:  track.Title,
		User:   fmUser,
	})
	if err == nil {
		playcount = info.UserPlaycount
		loved = info.UserLoved.Bool()
	}

	title := util.Possessive(fmUser)
	if track.NowPlaying {
		title += " Now Scrobbling"
	} else {
		title += " Last Scrobbled"
	}

	artistName := dctools.EscapeMarkdown(track.Artist.Name)
	trackName := dctools.EscapeMarkdown(track.Title)
	trackLink := dctools.Hyperlink(trackName, track.URL)

	trackField := fmt.Sprintf("%s - %s", artistName, trackLink)

	imageURL := track.Image.SizedURL(lastfm.ImgSizeLarge)
	if imageURL == "" {
		imageURL = lastfm.NoAlbumImageURL.Resize(lastfm.ImgSizeLarge)
	}

	footer := util.PluraliseWithCount("Scrobble", int64(playcount))
	if loved {
		footer = dctools.SeparateEmbedFooter(footer, "❤ Loved")
	}
	footer = dctools.SeparateEmbedFooter(footer, "Powered by Last.fm")

	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: title, URL: userURL(fmUser), Icon: scrobbleIcon,
		},
		Fields: []discord.EmbedField{
			{Name: "Track", Value: trackField},
		},
		Color:     scrobbleColour,
		Thumbnail: &discord.EmbedThumbnail{URL: imageURL},
		Footer:    &discord.EmbedFooter{Text: footer},
	}

	if track.Album.Title != "" {
		title := dctools.EscapeMarkdown(track.Album.Title)
		field := discord.EmbedField{Name: "Album", Value: title}
		embed.Fields = append(embed.Fields, field)
	}

	if !track.NowPlaying {
		embed.Timestamp = discord.Timestamp(track.ScrobbledAt)
	}

	ctx.RespondEmbed(embed)
}
