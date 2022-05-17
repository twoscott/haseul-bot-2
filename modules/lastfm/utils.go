package lastfm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/shkh/lastfm-go/lastfm"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

const (
	lastFmURL     = "https://www.last.fm"
	thumbURLFrame = "lastfm.freetls.fastly.net/i/u/174s/%s.jpg"
	imageURLFrame = "lastfm.freetls.fastly.net/i/u/%s.jpg"

	noArtistHash = "2a96cbd8b46e442fc41c2b86b821562f"
	noAlbumHash  = "c6f59c1e5e7240a4c0d427abd71f3dbb"
	noTrackHash  = "4128a6eb29f94943c9d206c08e625904"
	noAvatarHash = "818148bf682d429dc215c1705eb27b98"

	lastFmIcon   = "https://i.imgur.com/0OsgDI9.jpg"
	scrobbleIcon = "https://i.imgur.com/f89gJc7.jpg"
	artistIcon   = "https://i.imgur.com/FVvB8rx.jpg"
	albumIcon    = "https://i.imgur.com/9Srp3rE.jpg"
	trackIcon    = "https://i.imgur.com/4xUxZTk.jpg"

	lastFmColour   = 0xcf2b2b
	scrobbleColour = lastFmColour
	artistColour   = 0xf49d37
	albumColour    = 0x00ad7a
	trackColour    = 0x0066ff
)

func getThumbURL(hash string) string {
	return fmt.Sprintf(thumbURLFrame, hash)
}

func getImageURL(hash string) string {
	return fmt.Sprintf(imageURLFrame, hash)
}

func checkUserExists(user string) (bool, error) {
	if user == "" {
		err := errors.New("No Last.fm username provided")
		return false, err
	}

	userURL := "https://www.last.fm/user/" + user
	res, err := http.Head(userURL)
	if err != nil {
		return false, err
	}
	if res.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if res.StatusCode != 200 {
		err := errors.New(
			"Status code from Last.fm user exists check was neither 404 or 200",
		)
		return false, err
	}

	return true, nil
}

func getLfUser(ctx router.CommandCtx) (string, bool) {
	lfUser, err := db.LastFM.GetUser(ctx.Msg.Author.ID)
	if errors.Is(err, sql.ErrNoRows) {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please link a Last.fm username to your account using "+
				"`fm set`.",
		)
		return "", false
	}
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred trying to find your Last.fm username.",
		)
		return "", false
	}

	return lfUser, true
}

func getLfError(err error) *lastfm.LastfmError {
	lfErr, ok := err.(*lastfm.LastfmError)
	if !ok {
		lfErr = &lastfm.LastfmError{
			Code:    0,
			Message: "Unknown Last.fm Error",
		}
	}

	return lfErr
}

func getArtistLibraryURL(user string, tf timeframe) string {
	return lastFmURL + "/user/" + user +
		"/library/artists?date_preset=" + tf.datePreset
}

func getAlbumLibraryURL(user string, tf timeframe) string {
	return lastFmURL + "/user/" + user +
		"/library/albums?date_preset=" + tf.datePreset
}

func getTrackLibraryURL(user string, tf timeframe) string {
	return lastFmURL + "/user/" + user +
		"/library/tracks?date_preset=" + tf.datePreset
}

func getLibraryURL(user string) string {
	return lastFmURL + "/user/" + user + "/library"
}

func recentTracks(
	ctx router.CommandCtx,
	lfUser string,
	limit int) (*lastfm.UserGetRecentTracks, bool) {
	res, err := lf.User.GetRecentTracks(lastfm.P{
		"user": lfUser, "limit": limit,
	})
	if err != nil {
		lfErr := getLfError(err)
		switch lfErr.Code {
		case 6:
			dctools.ReplyWithWarning(ctx.State, ctx.Msg,
				fmt.Sprintf("Last.fm user %s does not exist.", lfUser),
			)
		case 8:
			log.Println(err)
			dctools.ReplyWithError(ctx.State, ctx.Msg,
				"I could not get a response from Last.fm. Please try again.",
			)
		default:
			log.Println(err)
			dctools.ReplyWithError(ctx.State, ctx.Msg,
				"Unknown Last.fm Error occurred.",
			)
		}
		return nil, false
	}

	if res.User == "" {
		res.User = lfUser
	}

	// correct extra track returned when now playing
	if len(res.Tracks) > limit {
		res.Tracks = res.Tracks[:limit]
	}

	return &res, true
}

func trackInfo(
	ctx router.CommandCtx,
	recentTracks *lastfm.UserGetRecentTracks) (*lastfm.TrackGetInfo, bool) {

	track1 := recentTracks.Tracks[0]
	lfUser := recentTracks.User

	res, err := lf.Track.GetInfo(lastfm.P{
		"artist":   track1.Artist.Name,
		"track":    track1.Name,
		"username": lfUser,
	})
	if err != nil {
		return nil, false
	}

	numPlayCount, err := strconv.Atoi(res.UserPlayCount)
	if err == nil {
		res.UserPlayCount = humanize.Comma(int64(numPlayCount))
	}

	return &res, true
}
