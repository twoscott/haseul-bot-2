package ytutil

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/twoscott/haseul-bot-2/utils/httputil"
)

const (
	youTubeURL = "https://www.youtube.com"
	queryPath  = "/results"
	videoOnly  = "EgIQAQ=="
)

var (
	ErrNoResultsFound  = errors.New("no results found")
	ErrUnknown         = errors.New("unknown error occurred")
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
func MultiSearch(query string, results int) ([]string, error) {
	return getResults(query, results)
}

func buildVideoURL(videoID string) string {
	return "https://youtu.be/" + videoID
}

func getResults(query string, resultLimit int) ([]string, error) {
	if query == "" {
		return nil, errors.New("no YouTube query provided")
	}

	queryURL, err := url.Parse(youTubeURL)
	if err != nil {
		return nil, err
	}

	queryURL.Path = queryPath
	queryBuilder := queryURL.Query()
	queryBuilder.Set("search_query", query)
	queryBuilder.Set("sp", videoOnly)
	queryURL.RawQuery = queryBuilder.Encode()

	res, err := httputil.Get(queryURL.String())
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
