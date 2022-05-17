package youtubeutil

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

const (
	youTubeQueryURL = "https://www.youtube.com/results"
	videoOnly       = "EgIQAQ=="
)

var (
	ErrNoResultsFound  = errors.New("No results found")
	ErrUnknown         = errors.New("Unknown error occurred")
	ytInitialDataRegex = regexp.MustCompile(`var ytInitialData\s*=.+?;`)
	ytContentRegex     = regexp.MustCompile(
		`"videoRenderer":\{"videoId":"(.+?)"`,
	)
)

// Search returns a single youtube video link for the given query.
func Search(query string) (string, error) {
	results, err := getResults(query, 1)
	if err != nil {
		return "", err
	}

	return results[0], err
}

// MultiSearch returns up to 20 youtube video links for the given query.
func MultiSearch(query string) ([]string, error) {
	return getResults(query, 20)
}

func buildVideoURL(videoID string) string {
	return "https://youtu.be/" + videoID
}

func getResults(query string, resultLimit int) ([]string, error) {
	queryURL, err := url.Parse(youTubeQueryURL)
	if err != nil {
		return nil, err
	}

	queryBuilder := queryURL.Query()
	queryBuilder.Set("search_query", query)
	queryBuilder.Set("sp", videoOnly)
	queryURL.RawQuery = queryBuilder.Encode()

	res, err := http.Get(queryURL.String())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err := errors.New(res.Status)
		log.Println("YouTube error:", err)
		return nil, err
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	ytDataBytes := ytInitialDataRegex.Find(bytes)
	if ytDataBytes == nil {
		return nil, ErrUnknown
	}

	renderersMatch := ytContentRegex.FindAllSubmatch(ytDataBytes, resultLimit)
	if renderersMatch == nil {
		return nil, ErrNoResultsFound
	}
	
	searchResults := make([]string, 0, resultLimit)
	for _, match := range renderersMatch {
		videoID := string(match[1])
		searchResults = append(searchResults, buildVideoURL(videoID))
	}

	return searchResults, nil
}
