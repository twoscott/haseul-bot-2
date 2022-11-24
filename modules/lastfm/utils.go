package lastfm

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/shkh/lastfm-go/lastfm"
	"github.com/twoscott/haseul-bot-2/utils/httputil"
)

const (
	lastFmURL     = "https://www.last.fm"
	thumbURLFrame = "https://lastfm.freetls.fastly.net/i/u/174s/%s"
	imageURLFrame = "https://lastfm.freetls.fastly.net/i/u/300x300/%s"

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

var imageRegexp = regexp.MustCompile(
	`https?://lastfm\.freetls\.fastly\.net/i/u/(.+?)/(.+?)(?:\.|$)`,
)

func toThumbnail(url string) string {
	match := imageRegexp.FindStringSubmatch(url)
	if match == nil {
		return ""
	}

	hash := match[2]
	return getThumbURL(hash)
}

func toImage(url string) string {
	match := imageRegexp.FindStringSubmatch(url)
	if match == nil {
		return ""
	}

	hash := match[2]
	return getImageURL(hash)
}

func getThumbURL(hash string) string {
	return fmt.Sprintf(thumbURLFrame, hash)
}

func getImageURL(hash string) string {
	return fmt.Sprintf(imageURLFrame, hash)
}

func checkUserExists(user string) (bool, error) {
	if user == "" {
		err := errors.New("no Last.fm username provided")
		return false, err
	}

	userURL := "https://www.last.fm/user/" + user
	res, err := httputil.Head(userURL)
	if err != nil {
		return false, err
	}
	if res.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if res.StatusCode != 200 {
		err := errors.New(
			"status code from Last.fm user exists check was neither 404 or 200",
		)
		return false, err
	}

	return true, nil
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
	return fmt.Sprintf(
		"%s/user/%s/library/artists?date_preset=%s",
		lastFmURL, user, tf.datePreset,
	)
}

func getAlbumLibraryURL(user string, tf timeframe) string {
	return fmt.Sprintf(
		"%s/user/%s/library/albums?date_preset=%s",
		lastFmURL, user, tf.datePreset,
	)
}

func getTrackLibraryURL(user string, tf timeframe) string {
	return fmt.Sprintf(
		"%s/user/%s/library/tracks?date_preset=%s",
		lastFmURL, user, tf.datePreset,
	)
}

func getLibraryURL(user string) string {
	return fmt.Sprintf("%s/user/%s/library", lastFmURL, user)
}

func getRecentTracks(
	user string, limit int64) (*lastfm.UserGetRecentTracks, error) {

	res, err := lf.User.GetRecentTracks(lastfm.P{
		"user": user, "limit": limit,
	})
	if err != nil {
		return nil, err
	}

	if res.User == "" {
		res.User = user
	}

	// correct extra track returned when a track is now playing.
	if int64(len(res.Tracks)) > limit {
		res.Tracks = res.Tracks[:limit]
	}

	return &res, nil
}

func getTrackInfo(
	recentTracks *lastfm.UserGetRecentTracks) (*lastfm.TrackGetInfo, error) {

	track1 := recentTracks.Tracks[0]
	lfUser := recentTracks.User

	res, err := lf.Track.GetInfo(lastfm.P{
		"artist":   track1.Artist.Name,
		"track":    track1.Name,
		"username": lfUser,
	})
	if err != nil {
		return nil, err
	}

	numPlayCount, err := strconv.Atoi(res.UserPlayCount)
	if err == nil {
		res.UserPlayCount = humanize.Comma(int64(numPlayCount))
	}

	return &res, nil
}

func errorResponseMessage(err error) string {
	lfErr := getLfError(err)
	switch lfErr.Code {
	case 6:
		return "Invalid parameters provided."
	case 8, 11, 16:
		return "Unable to get a response from Last.fm. Please try again."
	default:
		return "Unknown Last.fm Error occurred."
	}
}
