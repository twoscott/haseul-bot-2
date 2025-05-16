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

var lastFMRecentCommand = &router.SubCommand{
	Name:        "recents",
	Description: "Displays your recently scrobbled tracks",
	Handler: &router.CommandHandler{
		Executor: lastFMRecentExec,
		Defer:    true,
	},
	Options: []discord.CommandOptionValue{
		&discord.IntegerOption{
			OptionName:  "tracks",
			Description: "The number of recent tracks to display for the user",
			Min:         option.NewInt(1),
			Max:         option.NewInt(1000),
		},
	},
}

func lastFMRecentExec(ctx router.CommandCtx) {
	trackCount, _ := ctx.Options.Find("tracks").IntValue()
	if trackCount == 0 {
		trackCount = 10
	}

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

	res, err := fm.User.RecentTracks(lastfm.RecentTracksParams{
		User:  fmUser,
		Limit: uint(trackCount),
	})
	if err != nil {
		log.Println(err)
		if msg, ok := errMessage(err); ok {
			ctx.RespondError(msg)
		} else {
			ctx.RespondError("Unable to fetch your recent scrobbles from Last.fm.")
		}
		return
	}

	if len(res.Tracks) < 1 {
		ctx.RespondWarning("You have not scrobbled any tracks on Last.fm.")
		return
	}

	if len(res.Tracks) > int(trackCount) {
		res.Tracks = res.Tracks[:trackCount]
	}

	trackList := make([]string, 0, len(res.Tracks))
	for _, track := range res.Tracks {
		artistName := dctools.EscapeMarkdown(track.Artist.Name)
		trackName := dctools.EscapeMarkdown(track.Title)
		trackLink := dctools.Hyperlink(trackName, track.URL)

		line := fmt.Sprintf("- %s - %s", artistName, trackLink)

		if track.NowPlaying {
			line += " (Now Playing)"
		} else {
			t := track.ScrobbledAt.Time()
			line += " " + dctools.TimestampStyled(t, dctools.RelativeTime)
		}

		trackList = append(trackList, line)
	}

	pages := util.PagedLines(trackList, 2048, 25)

	title := util.Possessive(fmUser) + " Recent Scrobbles"

	imageURL := res.Tracks[0].Image.SizedURL(lastfm.ImgSizeLarge)
	if imageURL == "" {
		imageURL = lastfm.NoAlbumImageURL.Resize(lastfm.ImgSizeLarge)
	}

	footer := util.PluraliseWithCount("Total Scrobble", int64(res.Total))
	footer = dctools.SeparateEmbedFooter(footer, "Powered by Last.fm")

	messagePages := make([]router.MessagePage, len(pages))
	for i, page := range pages {
		id := fmt.Sprintf("Page %d/%d", i+1, len(pages))

		e := discord.Embed{
			Author: &discord.EmbedAuthor{
				Name: title, URL: libraryURL(fmUser), Icon: scrobbleIcon,
			},
			Description: page,
			Color:       scrobbleColour,
			Thumbnail:   &discord.EmbedThumbnail{URL: imageURL},
			Footer: &discord.EmbedFooter{
				Text: dctools.SeparateEmbedFooter(id, footer),
			},
		}

		messagePages[i] = router.MessagePage{Embeds: []discord.Embed{e}}
	}

	ctx.RespondPaging(messagePages)
}
