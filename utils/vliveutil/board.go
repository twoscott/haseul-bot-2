package vliveutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Board struct {
	ID          int64  `json:"boardId"`
	Title       string `json:"title"`
	Type        string `json:"boardType"`
	ChannelCode string `json:"channelCode"`
}

const boardFields = "boardId,title,boardType,channelCode"

// BoardExists returns whether or not a VLIVE board with the provided board
// ID exists.
func BoardExists(boardID int64) (bool, error) {
	boardURL, err := buildBoardURL(boardID)
	if err != nil {
		return false, err
	}

	res, err := headRequest(*boardURL)
	if err != nil {
		return false, err
	}

	return res.StatusCode == http.StatusOK, nil
}

// GetBoard returns the VLIVE board with the corresponding provided board ID.
func GetBoard(boardID int64) (*Board, *http.Response, error) {
	boardURL, err := buildBoardURL(boardID)
	if err != nil {
		return nil, nil, err
	}

	bytes, res, err := getRequestBytes(*boardURL)
	if err != nil {
		return nil, nil, err
	}

	var board Board
	json.Unmarshal(bytes, &board)

	return &board, res, nil
}

func buildBoardURL(boardID int64) (*url.URL, error) {
	endpoint := boardPath(boardID)

	queryBuilder := vliveQueryBuilder()
	queryBuilder.Set("fields", boardFields)

	channelURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	channelURL.RawQuery = queryBuilder.Encode()
	return channelURL, nil
}

func boardPath(boardID int64) string {
	return BoardEndpoint + fmt.Sprintf("/board-%d", boardID)
}
