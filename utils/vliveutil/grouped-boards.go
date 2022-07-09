package vliveutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type BoardGroup struct {
	GroupTitle string  `json:"groupTitle"`
	Boards     []Board `json:"boards"`
}

// GetGroupedBoards returns all of a VLIVE channel's board groups and their
// contained boards.
func GetGroupedBoards(
	channelCode string) ([]BoardGroup, *http.Response, error) {

	boardsURL, err := buildGroupedBoardsURL(channelCode)
	if err != nil {
		return nil, nil, err
	}

	bytes, res, err := getRequestBytes(*boardsURL)
	if err != nil {
		return nil, nil, err
	}

	var groups []BoardGroup
	json.Unmarshal(bytes, &groups)

	return groups, res, nil
}

// GetUnwrappedBoards returns all of a VLIVE channel's boards directly,
// ignoring the encasing groups.
func GetUnwrappedBoards(channelCode string) ([]Board, *http.Response, error) {
	groups, res, err := GetGroupedBoards(channelCode)
	if err != nil {
		return nil, res, err
	}

	boards := make([]Board, 0, len(groups))
	for _, g := range groups {
		for _, b := range g.Boards {
			boards = append(boards, b)
		}
	}

	return boards, res, err
}

func buildGroupedBoardsURL(channelCode string) (*url.URL, error) {
	endpoint := groupedBoardsPath(channelCode)

	queryBuilder := vliveQueryBuilder()
	queryBuilder.Set("fields", boardFields)

	channelURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	channelURL.RawQuery = queryBuilder.Encode()
	return channelURL, nil
}

func groupedBoardsPath(channelCode string) string {
	return BoardEndpoint + fmt.Sprintf("/channel-%s/groupedBoards", channelCode)
}
