package vliveutil

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Channel struct {
	ProfileImage  string `json:"profileImg"`
	Name          string `json:"name"`
	LastUpdatedAt int64  `json:"latestUpdatedAt"`
	Code          string `json:"channelCode"`
}

// ChannelExists returns whether or not a VLIVE channel with the provided
// channel code exists.
func ChannelExists(channelCode string) (bool, error) {
	channelURL, err := buildChannelURL(channelCode)
	if err != nil {
		return false, err
	}

	res, err := headRequest(*channelURL)
	if err != nil {
		return false, err
	}

	return res.StatusCode == http.StatusOK, nil
}

// GetChannel returns the VLIVE channel with the provided channel code.
func GetChannel(channelCode string) (*Channel, *http.Response, error) {
	channelURL, err := buildChannelURL(channelCode)
	if err != nil {
		return nil, nil, err
	}

	bytes, res, err := getRequestBytes(*channelURL)
	if err != nil {
		return nil, nil, err
	}

	var channel Channel
	json.Unmarshal(bytes, &channel)

	return &channel, res, nil
}

func buildChannelURL(channelCode string) (*url.URL, error) {
	endpoint := ChannelsEndpoint + "/" + channelCode

	queryBuilder := vliveQueryBuilder()

	channelURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	channelURL.RawQuery = queryBuilder.Encode()
	return channelURL, nil
}
