package ytutil

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/twoscott/haseul-bot-2/utils/httputil"
	"golang.org/x/net/html"
)

const NumOfSuggestions = 10

const (
	youtubeSuggestionsURL    = "https://suggestqueries-clients6.youtube.com"
	autocompletePath         = "/complete/search"
	autocompleteClient       = "youtube"
	autocompleteHostLanguage = "en"
	autocompleteGeoLocation  = "us"
	autocompleteDS           = "yt"
)

var ytSuggestionsRegex = regexp.MustCompile(`\[\s*"(.+?)"`)

// GetSuggestions gets live search suggestions for a given YouTube query.
func GetSuggestions(query string) ([]string, error) {
	if query == "" {
		return nil, errors.New(
			"No YouTube query provided to get suggestions for",
		)
	}

	suggestURL, err := url.Parse(youtubeSuggestionsURL)
	if err != nil {
		return nil, err
	}

	suggestURL.Path = autocompletePath
	queryBuilder := suggestURL.Query()
	queryBuilder.Set("client", autocompleteClient)
	queryBuilder.Set("hl", autocompleteHostLanguage)
	queryBuilder.Set("gl", autocompleteGeoLocation)
	queryBuilder.Set("ds", autocompleteDS)
	queryBuilder.Set("q", query)
	suggestURL.RawQuery = queryBuilder.Encode()

	res, err := httputil.Get(suggestURL.String())
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

	return parseSuggestions(bytes)
}

// parseSuggestions fetches string sections from the raw data returned from
// YouTube and unquotes/unescapes the characters.
func parseSuggestions(bytes []byte) ([]string, error) {
	matches := ytSuggestionsRegex.FindAllStringSubmatch(
		string(bytes),
		NumOfSuggestions+1,
	)
	if matches == nil {
		return nil, errors.New("Unable to parse suggestions")
	}

	suggestions := make([]string, 0, len(matches))
	for _, match := range matches {
		suggestion, err := strconv.Unquote(`"` + match[1] + `"`)
		if err != nil {
			suggestion = match[1]
		}
		suggestion = html.UnescapeString(suggestion)
		suggestions = append(suggestions, suggestion)
	}

	// [0] contains the query provided - sometimes the first suggestions is
	// the same as the query so we should slice the first element off.
	if len(suggestions) >= 2 && suggestions[0] == suggestions[1] {
		suggestions = suggestions[1:]
	}

	// first string is the entered search query.
	return suggestions, nil
}
