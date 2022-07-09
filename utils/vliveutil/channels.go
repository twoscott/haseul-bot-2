package vliveutil

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

type ChannelResult struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
	Type string `json:"type"`
	Code string `json:"code"`
}

// SearchChannels returns channel results that are returned when queried with
// the provided query by their channel names.
func SearchChannels(
	query string, max uint) ([]ChannelResult, *http.Response, error) {

	channelsURL, err := buildChannelsURL(query, max)
	if err != nil {
		return nil, nil, err
	}

	bytes, res, err := getRequestBytes(*channelsURL)
	if err != nil {
		return nil, nil, err
	}

	var channels []ChannelResult
	json.Unmarshal(bytes, &channels)

	return channels, res, nil
}

func buildChannelsURL(query string, max uint) (*url.URL, error) {
	maxStr := strconv.FormatUint(uint64(max), 10)

	queryBuilder := vliveQueryBuilder()
	queryBuilder.Set("query", query)
	queryBuilder.Set("maxNumOfRows", maxStr)

	channelURL, err := url.Parse(SearchChannelsEndpoint)
	if err != nil {
		return nil, err
	}

	channelURL.RawQuery = queryBuilder.Encode()
	return channelURL, nil
}
