package lastfm

import (
	"errors"
	"fmt"

	"github.com/twoscott/gobble-fm/api"
	"github.com/twoscott/gobble-fm/lastfm"
)

const (
	lastFMIcon   = "https://i.imgur.com/0OsgDI9.jpg"
	scrobbleIcon = "https://i.imgur.com/f89gJc7.jpg"
	artistIcon   = "https://i.imgur.com/FVvB8rx.jpg"
	albumIcon    = "https://i.imgur.com/9Srp3rE.jpg"
	trackIcon    = "https://i.imgur.com/4xUxZTk.jpg"

	lastFMColour   = 0xcf2b2b
	scrobbleColour = lastFMColour
	artistColour   = 0xf49d37
	albumColour    = 0x00ad7a
	trackColour    = 0x0066ff
)

func errMessage(err error) (string, bool) {
	fmErr, ok := fmError(err)
	if !ok {
		return "", false
	}

	switch fmErr.Code {
	case api.ErrInvalidParameters:
		return "Invalid parameters provided.", true
	case api.ErrServiceOffline:
		return "Last.fm service is offline. Please try again later.", true
	case api.ErrOperationFailed, api.ErrServiceUnavailable:
		return "Unable to get a response from Last.fm. Please try again.", true
	default:
		return "Unknown Last.fm error occurred.", true
	}
}

func fmError(err error) (*api.LastFMError, bool) {
	var fmerr *api.LastFMError
	if errors.As(err, &fmerr) {
		return fmerr, true
	}

	return nil, false
}

func artistLibraryURL(user string, tf timeframe) string {
	return fmt.Sprintf(
		"%s/user/%s/library/artists?date_preset=%s",
		lastfm.BaseURL, user, tf.datePreset,
	)
}

func albumLibraryURL(user string, tf timeframe) string {
	return fmt.Sprintf(
		"%s/user/%s/library/albums?date_preset=%s",
		lastfm.BaseURL, user, tf.datePreset,
	)
}

func trackLibraryURL(user string, tf timeframe) string {
	return fmt.Sprintf(
		"%s/user/%s/library/tracks?date_preset=%s",
		lastfm.BaseURL, user, tf.datePreset,
	)
}

func libraryURL(user string) string {
	return fmt.Sprintf("%s/user/%s/library", lastfm.BaseURL, user)
}

func userURL(user string) string {
	return fmt.Sprintf("%s/user/%s", lastfm.BaseURL, user)
}
