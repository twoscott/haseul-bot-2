package lastfm

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/twoscott/haseul-bot-2/utils/httputil"
)

var errScrapingArtistImage = errors.New("unable to scrape Last.fm artist image")

func scrapeArtistImage(artist string) (string, error) {
	artistImagesURL, err := url.Parse(lastFmURL)
	if err != nil {
		return "", err
	}

	artistNamePath := url.PathEscape(artist)
	artistImagesURL.Path = "music/" + artistNamePath + "/+images"

	res, err := httputil.Get(artistImagesURL.String())
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

	imageItem := doc.Find(".image-list-item").First()
	image := imageItem.ChildrenFiltered("img").First()
	imageURL := image.AttrOr("src", getThumbURL(noArtistHash))
	imageURL = toImage(imageURL)

	return imageURL, nil
}
