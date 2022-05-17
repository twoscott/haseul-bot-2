package util

import (
	"net/http"
	"strings"
	"time"
)

// ImgUploadTime returns the time an image was uploaded
// (or its last-modified header)
func ImgUploadTime(url string) (time.Time, error) {
	res, err := http.Head(url)
	if err != nil {
		return time.Time{}, err
	}

	return HeaderModifiedTime(res.Header)
}

// HeaderModifiedTime returns a time object representing when the object of the
// header was last modified.
func HeaderModifiedTime(header http.Header) (time.Time, error) {
	modified := header.Get("Last-Modified")
	if modified == "" {
		return time.Time{}, nil
	}

	uploadTime, err := time.Parse(time.RFC1123, modified)
	return uploadTime, err
}

// PrettifyContentType returns the latter section of the content type; after
// the "/", converted to uppercase.
func PrettifyContentType(contentType string) string {
	if contentType == "" {
		return "Unknown"
	}

	parts := strings.Split(contentType, "/")
	if len(parts) < 2 {
		return strings.ToUpper(contentType)
	}

	return strings.ToUpper(parts[1])
}
