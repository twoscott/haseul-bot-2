package vliveutil

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type BoardPostsPager struct {
	PostTimestamp int64
	PostID        string
}

func NewBoardPostsPager(postTimestamp int64, postID string) *BoardPostsPager {
	return &BoardPostsPager{
		PostTimestamp: postTimestamp,
		PostID:        postID,
	}
}

func (p BoardPostsPager) String() string {
	return fmt.Sprintf("%d,%s", p.PostTimestamp, p.PostID)
}

// GetBoardPosts returns limit number of recent posts from a channel's VLIVE
// board at board ID.
func GetBoardPosts(boardID int64, limit int) ([]Post, *http.Response, error) {
	postsURL, err := buildBoardPostsURL(boardID, limit, nil)
	if err != nil {
		return nil, nil, err
	}

	posts, res, err := getPosts(*postsURL)
	if err != nil {
		return nil, nil, err
	}

	return posts, res, nil
}

// GetBoardPosts returns limit number of recent posts from a channel's VLIVE
// board at board ID that were posted before the provided before parameter.
func GetBoardPostsBefore(
	boardID int64,
	limit int,
	before BoardPostsPager) ([]Post, *http.Response, error) {

	postsURL, err := buildBoardPostsURL(boardID, limit, &before)
	if err != nil {
		return nil, nil, err
	}

	posts, res, err := getPosts(*postsURL)
	if err != nil {
		return nil, nil, err
	}

	return posts, res, nil
}

func buildBoardPostsURL(
	boardID int64, limit int, before *BoardPostsPager) (*url.URL, error) {

	endpoint := boardPostsPath(boardID)

	queryBuilder := vliveQueryBuilder()
	queryBuilder.Set("fields", postFields)
	queryBuilder.Set("sortType", "LATEST")

	if limit > 0 {
		queryBuilder.Set("limit", strconv.Itoa(limit))
	}

	if before != nil {
		queryBuilder.Set("before", before.String())
	}

	postsURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	postsURL.RawQuery = queryBuilder.Encode()
	return postsURL, nil
}

func boardPostsPath(boardID int64) string {
	return PostEndpoint + fmt.Sprintf("/board-%d/posts", boardID)
}
