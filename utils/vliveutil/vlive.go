package vliveutil

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/twoscott/haseul-bot-2/utils/httputil"
)

const (
	BaseURL = "https://www.vlive.tv"

	PostPath     = "/globalv-web/vam-web/post/v1.0"
	PostEndpoint = BaseURL + PostPath

	BoardPath     = "/globalv-web/vam-web/board/v1.0"
	BoardEndpoint = BaseURL + BoardPath

	ChannelsPath     = "/globalv-web/vam-web/vhs/store/v1.0/channels"
	ChannelsEndpoint = BaseURL + ChannelsPath

	SearchChannelsPath     = "/search/auto/channels"
	SearchChannelsEndpoint = BaseURL + SearchChannelsPath

	AppID        = "8c6cc7b45d2568fb668be6e05b6e5a3b"
	PlatformType = "ANDROID"
	GeoLocation  = "KR"
	Locale       = "en_US"
)

var vliveClient = httputil.NewClient(10 * time.Second)

func newGetRequest(requestURL url.URL) *http.Request {
	return newRequest(requestURL, "GET")
}

func newHeadRequest(requestURL url.URL) *http.Request {
	return newRequest(requestURL, "HEAD")
}

func newRequest(requestURL url.URL, method string) *http.Request {
	header := http.Header{}
	header.Set("Referer", BaseURL)

	return &http.Request{
		Method: method,
		URL:    &requestURL,
		Header: header,
	}
}

func getRequest(requestURL url.URL) (*http.Response, error) {
	req := newGetRequest(requestURL)
	return vliveClient.Do(req)
}

func getRequestBytes(requestURL url.URL) ([]byte, *http.Response, error) {
	res, err := getRequest(requestURL)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res, err
	}

	return b, res, nil
}

func headRequest(requestURL url.URL) (*http.Response, error) {
	req := newHeadRequest(requestURL)
	return vliveClient.Do(req)
}

func headRequestBytes(requestURL url.URL) ([]byte, *http.Response, error) {
	res, err := headRequest(requestURL)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res, err
	}

	return b, res, nil
}

func vliveQueryBuilder() url.Values {
	queryBuilder := url.Values{}
	queryBuilder.Set("appId", AppID)
	queryBuilder.Set("platformType", PlatformType)
	queryBuilder.Set("gcc", GeoLocation)
	queryBuilder.Set("locale", Locale)

	return queryBuilder
}
