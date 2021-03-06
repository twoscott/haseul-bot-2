package lastfm

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/twoscott/haseul-bot-2/utils/httputil"
)

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

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	imageItem := doc.Find(".image-list-item").First()
	image := imageItem.ChildrenFiltered("img").First()
	imageURL := image.AttrOr("src", getThumbURL(noArtistHash))
	imageURL = strings.Replace(imageURL, "/i/u/avatar170s", "/i/u/300x300", 1)

	return imageURL, nil
}

// func scrapeTopArtistImages(...)
