package httputil

import (
	"io"
	"net/http"
	"time"
)

var defaultClient = http.Client{
	Timeout: 5 * time.Second,
}

// NewHttpClient provides a clean way to create a new HTTP client with a
// custom request timeout duration.
func NewClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

// Do sends a HTTP request with a 5 second timeout.
func Do(req *http.Request) (*http.Response, error) {
	return defaultClient.Do(req)
}

// Get sends a HTTP GET request with a 5 second timeout.
func Get(url string) (*http.Response, error) {
	return defaultClient.Get(url)
}

// Head sends a HTTP HEAD request with a 5 second timeout.
func Head(url string) (*http.Response, error) {
	return defaultClient.Head(url)
}

// Post sends a HTTP POST request with a 5 second timeout.
func Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return defaultClient.Post(url, contentType, body)
}

// ImgUploadTime returns the time an image was uploaded
// (or its last-modified header)
func ImgUploadTime(url string) (time.Time, error) {
	res, err := http.Head(url)
	if err != nil {
		return time.Time{}, err
	}

	return HeaderModifiedTime(res.Header)
}
