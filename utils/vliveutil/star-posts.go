package vliveutil

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GetStarPosts returns limit number of recent posts from a channelCode's
// VLIVE star.
func GetStarPosts(
	channelCode string, limit int) ([]Post, *http.Response, error) {

	postsURL, err := buildStarPostsURL(channelCode, limit)
	if err != nil {
		return nil, nil, err
	}

	posts, res, err := getPosts(*postsURL)
	if err != nil {
		return nil, nil, err
	}

	return posts, res, nil
}

func buildStarPostsURL(channelCode string, limit int) (*url.URL, error) {
	endpoint := starPostsPath(channelCode)

	queryBuilder := vliveQueryBuilder()
	queryBuilder.Set("limit", strconv.Itoa(limit))
	queryBuilder.Set("fields", postFields)

	postsURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	postsURL.RawQuery = queryBuilder.Encode()
	return postsURL, nil
}

func starPostsPath(channelCode string) string {
	return PostEndpoint + fmt.Sprintf("/channel-%s/starPosts", channelCode)
}
