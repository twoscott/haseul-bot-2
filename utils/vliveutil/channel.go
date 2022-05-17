package vliveutil

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Channel struct {
	ProfileImage  string `json:"profileImg"`
	Name          string `json:"name"`
	LastUpdatedAt int64  `json:"latestUpdatedAt"`
	Code          string `json:"channelCode"`
}

func ChannelExists(channelCode string) (bool, error) {
	channelURL, err := BuildChannelURL(channelCode)
	if err != nil {
		return false, err
	}

	req := NewGetRequest(channelURL)
	res, err := vliveClient.Do(req)
	if err != nil {
		return false, err
	}

	return res.StatusCode == http.StatusOK, nil
}

func GetChannel(channelCode string) (*Channel, error) {
	channelURL, err := BuildChannelURL(channelCode)
	if err != nil {
		return nil, err
	}

	req := NewGetRequest(channelURL)
	res, err := vliveClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var channel Channel
	json.Unmarshal(bytes, &channel)

	return &channel, nil
}

func BuildChannelURL(channelCode string) (*url.URL, error) {
	endpoint := ChannelEndpoint + "/" + channelCode

	queryBuilder := url.Values{}
	queryBuilder.Set("appId", AppID)
	queryBuilder.Set("platformType", "PC")

	channelURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	channelURL.RawQuery = queryBuilder.Encode()
	return channelURL, nil
}
