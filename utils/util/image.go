package util

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net/http"

	"github.com/twoscott/haseul-bot-2/utils/httputil"
)

type ImageType uint8

const (
	Unknown ImageType = iota
	JPEG
	PNG
	GIF
)

// String returns the string representation of the image type.
func (t ImageType) String() string {
	switch t {
	case JPEG:
		return "JPEG"
	case PNG:
		return "PNG"
	case GIF:
		return "GIF"
	default:
		return "Unknown"
	}
}

// Image represents a raw image.
type Image struct {
	data []byte
}

// NewImage returns a new instance of image containing data.
func NewImage(data []byte) *Image {
	return &Image{data: data}
}

func ImageFromURL(url string) (*Image, *http.Response, error) {
	res, err := httputil.Get(url)
	if err != nil {
		return nil, res, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, res, errors.New("Bad response")
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res, err
	}

	return &Image{data: bytes}, res, nil
}

// Size returns the image size in bytes.
func (i Image) Size() int {
	return len(i.data)
}

// Dimensions returns the width and height of the image.
func (i Image) Dimensions() [2]uint32 {
	switch i.Type() {
	case JPEG:
		SOF0 := bytes.Index(i.data, []byte{0xFF, 0xC0})
		width := binary.BigEndian.Uint16(i.data[SOF0+7 : SOF0+9])
		height := binary.BigEndian.Uint16(i.data[SOF0+5 : SOF0+7])
		return [2]uint32{uint32(width), uint32(height)}
	case PNG:
		width := binary.BigEndian.Uint32(i.data[16:20])
		height := binary.BigEndian.Uint32(i.data[20:24])
		return [2]uint32{width, height}
	case GIF:
		width := binary.LittleEndian.Uint16(i.data[6:8])
		height := binary.LittleEndian.Uint16(i.data[8:10])
		return [2]uint32{uint32(width), uint32(height)}
	default:
		return [2]uint32{0, 0}
	}
}

// Width returns the width of the image.
func (i Image) Width() uint32 {
	return i.Dimensions()[0]
}

// Height returns the height of the image.
func (i Image) Height() uint32 {
	return i.Dimensions()[1]
}

// Type returns the image's type.
func (i Image) Type() ImageType {
	switch {
	case i.IsJPEG():
		return JPEG
	case i.IsPNG():
		return PNG
	case i.IsGIF():
		return GIF
	default:
		return Unknown
	}
}

// IsJPEG returns whether the image is a JPEG.
func (i Image) IsJPEG() bool {
	header := i.data[:2]
	return bytes.Equal(header, []byte{0xFF, 0xD8})
}

// IsPNG returns whether the image is a PNG.
func (i Image) IsPNG() bool {
	header := i.data[:4]
	return bytes.Equal(header, []byte{0x89, 'P', 'N', 'G'})
}

// IsGIF returns whether the image is a GIF.
func (i Image) IsGIF() bool {
	header := i.data[:3]
	return bytes.Equal(header, []byte{'G', 'I', 'F'})
}
