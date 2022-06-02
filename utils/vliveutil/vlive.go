package vliveutil

import (
	"net/http"
	"net/url"
	"time"

	"github.com/twoscott/haseul-bot-2/utils/httputil"
)

const (
	BaseURL = "https://www.vlive.tv"

	PostPath     = "/globalv-web/vam-web/post/v1.0"
	PostEndpoint = BaseURL + PostPath

	ChannelPath     = "/globalv-web/vam-web/vhs/store/v1.0/channels"
	ChannelEndpoint = BaseURL + ChannelPath

	AppID = "8c6cc7b45d2568fb668be6e05b6e5a3b"
)

var vliveClient = httputil.NewClient(10 * time.Second)

func NewGetRequest(requestURL *url.URL) *http.Request {
	return newRequest(requestURL, "GET")
}

func NewHeadRequest(requestURL *url.URL) *http.Request {
	return newRequest(requestURL, "HEAD")
}

func newRequest(requestURL *url.URL, method string) *http.Request {
	header := http.Header{}
	header.Set("Referer", BaseURL)

	return &http.Request{
		Method: method,
		URL:    requestURL,
		Header: header,
	}
}
