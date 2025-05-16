package lastfm

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/twoscott/gobble-fm/lastfm"
	"github.com/twoscott/haseul-bot-2/utils/httputil"
)

var errScrapingArtistImage = errors.New("unable to scrape Last.fm artist image")

func scrapeArtistImage(artist string, size lastfm.ImgSize) (string, error) {
	artistURL, err := url.Parse(lastfm.BaseURL)
	if err != nil {
		return "", err
	}

	artistURL.Path = "music/" + url.PathEscape(artist)

	res, err := httputil.Get(artistURL.String())
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", errScrapingArtistImage
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	noArtistImg := lastfm.NoArtistImageURL.Resize(size)
	image := doc.Find(`meta[property="og:image"]`).First()
	imageURL := image.AttrOr("content", noArtistImg)
	imageURL = lastfm.ImageURL(imageURL).Resize(size)

	return imageURL, nil
}
