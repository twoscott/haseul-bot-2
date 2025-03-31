package util

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"io"
	"net/http"

	colx "github.com/marekm4/color-extractor"
	"github.com/twoscott/haseul-bot-2/utils/httputil"

	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// GetImageURL returns a HTTP response from requesting an image URL.
func GetImageURL(url string) (*http.Response, error) {
	res, err := httputil.Get(url)
	if err != nil {
		return res, err
	}

	if res.StatusCode != http.StatusOK {
		return res, errors.New("bad response")
	}

	return res, nil
}

// ImageFromResponse converts a HTTP response containing an image body to a Go
// image.
func ImageFromResponse(res http.Response) (image.Image, error) {
	img, _, err := image.Decode(res.Body)
	return img, err
}

// DataFromImageURL returns the raw data from an image URL.
func DataFromImageURL(url string) ([]byte, error) {
	res, err := GetImageURL(url)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// ImageFromURL returns a Go image interface representation of an image at URL.
func ImageFromURL(url string) (image.Image, error) {
	res, err := GetImageURL(url)
	if err != nil {
		return nil, err
	}

	return ImageFromResponse(*res)
}

// ColourFromImage returns a set of colours extracted from the given image.
func ColourFromImage(img image.Image) color.Color {
	c := colx.ExtractColors(img)
	if len(c) < 1 {
		return color.Black
	}

	return c[0]
}

// ColourFromURL fetches the image at URL and decodes it into an image, then
// extracts colours from that image.
func ColourFromURL(url string) (color.Color, error) {
	img, err := ImageFromURL(url)
	if err != nil {
		return nil, err
	}

	return ColourFromImage(img), nil
}

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

// RawImage represents a raw image.
type RawImage struct {
	data []byte
}

// NewRawImage returns a new instance of RawImage containing data.
func NewRawImage(data []byte) *RawImage {
	return &RawImage{data: data}
}

// RawImageFromURL returns a raw representation of an image with useful helper
// methods for finding the size, dimensions, and image type.
func RawImageFromURL(url string) (*RawImage, error) {
	data, err := DataFromImageURL(url)
	if err != nil {
		return nil, err
	}

	return NewRawImage(data), nil
}

// RawImageFromResponse converts a HTTP response containing an image body to a
// RawImage.
func RawImageFromResponse(res http.Response) (*RawImage, error) {
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return NewRawImage(bytes), err
}

// ToImage converts RawImage to a Go Image and returns it.
func (i RawImage) ToImage() (image.Image, error) {
	r := bytes.NewReader(i.data)

	img, _, err := image.Decode(r)
	return img, err
}

// Colour extracts a prominent colour from the image and returns it.
func (i RawImage) Colour() (color.Color, error) {
	c, err := i.ToImage()
	if err != nil {
		return nil, err
	}

	return ColourFromImage(c), nil
}

// Size returns the image size in bytes.
func (i RawImage) Size() int {
	return len(i.data)
}

// Dimensions returns the width and height of the image.
func (i RawImage) Dimensions() [2]uint32 {
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
func (i RawImage) Width() uint32 {
	return i.Dimensions()[0]
}

// Height returns the height of the image.
func (i RawImage) Height() uint32 {
	return i.Dimensions()[1]
}

// Type returns the image's type.
func (i RawImage) Type() ImageType {
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
func (i RawImage) IsJPEG() bool {
	header := i.data[:2]
	return bytes.Equal(header, []byte{0xFF, 0xD8})
}

// IsPNG returns whether the image is a PNG.
func (i RawImage) IsPNG() bool {
	header := i.data[:4]
	return bytes.Equal(header, []byte{0x89, 'P', 'N', 'G'})
}

// IsGIF returns whether the image is a GIF.
func (i RawImage) IsGIF() bool {
	header := i.data[:3]
	return bytes.Equal(header, []byte{'G', 'I', 'F'})
}
