package dctools

import (
	"errors"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

const (
	BlurpleColour = 0x5865f2
	GreenColour   = 0x57f287
	YellowColour  = 0xfee75c
	FuchsiaColour = 0xeb459e
	RedColour     = 0xed4245
	WhiteColour   = 0xffffff
	BlackColour   = 0x23272a
	// EmbedBackColour is the default embed background colour Discord uses.
	EmbedBackColour = 0x000000
)

var ErrInvalidHexColour = errors.New("invalid hexadecimal RGB value provided")

// ColourInvalid returns whether a colour is either null or a default colour.
func ColourInvalid(colour discord.Color) bool {
	return colour == discord.NullColor || colour == 0x000000
}

// EmbedColour returns a Discord colour where default black colours are
// replaced with the default embed background colour.
func EmbedColour(colour discord.Color) discord.Color {
	if ColourInvalid(colour) {
		colour = 0x000000
	}
	return colour
}

// EmbedImageColour returns an embed colour for a Discord image URL.
func EmbedImageColour(url string) (discord.Color, error) {
	c, err := util.ColourFromURL(ResizeImage(url, 256))
	if err != nil {
		return discord.NullColor, err
	}

	return ConvertColour(c), nil
}

// HexToColour converts a string representing a hexadecimal colour value
// to a Discord integer colour.
// e.g: "#3251cf" > 0x3251cf
//
// String must either begin with a '#' or consist only of hexadecimal values.
func HexToColour(hex string) (discord.Color, error) {
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) > 6 {
		return discord.NullColor, ErrInvalidHexColour
	}

	colourInt64, err := strconv.ParseInt(hex, 16, 32)
	if err != nil {
		return discord.NullColor, err
	}

	return discord.Color(colourInt64), nil
}

// ConverColor converts a Go colour into a Discord colour.
func ConvertColour(c color.Color) discord.Color {
	return RGBAToColour(c.RGBA())
}

// HSVToColour takes 3 values and converts it to a Discord colour;
// - hue - 0-360
// - sat - 0-1
// - val - 0-1
func HSVToColour(hue, sat, val float64) discord.Color {
	c := val * sat
	arc := hue / 60
	x := c * (1 - math.Abs(math.Mod(arc, 2)-1))

	var rf, gf, bf float64

	switch {
	case arc < 1:
		rf, gf, bf = c, x, 0
	case arc < 2:
		rf, gf, bf = x, c, 0
	case arc < 3:
		rf, gf, bf = 0, c, x
	case arc < 4:
		rf, gf, bf = 0, x, c
	case arc < 5:
		rf, gf, bf = x, 0, c
	case arc < 6:
		rf, gf, bf = c, 0, x
	}

	m := val - c
	r := discord.Color(math.Round((rf + m) * 255))
	g := discord.Color(math.Round((gf + m) * 255))
	b := discord.Color(math.Round((bf + m) * 255))

	return r<<16 | g<<8 | b
}

// RGBAtoRGB takes 4 RGBA (16-bit premultiplied) values and converts it to
// 8-bit RGB (0-255);
// - red   - 0-0xffff
// - green - 0-0xffff
// - blue  - 0-0xffff
// - alpha - 0-0xffff
func RGBAToRGB(red, green, blue, alpha uint32) (uint8, uint8, uint8) {
	if alpha == 0 {
		return 0, 0, 0 // Avoid division by zero
	}

	// Undo alpha-premultiplication (c = original * a / 0xffff)
	red = (red * 0xffff) / alpha
	green = (green * 0xffff) / alpha
	blue = (blue * 0xffff) / alpha

	// Scale from 16-bit to 8-bit
	r8 := uint8(red >> 8)
	g8 := uint8(green >> 8)
	b8 := uint8(blue >> 8)

	return r8, g8, b8
}

// RGBToColour takes 4 RGBA values and converts it to a Discord colour;
// - red   - 0-0xffff
// - green - 0-0xffff
// - blue  - 0-0xffff
// - alpha - 0-0xffff
func RGBAToColour(red, green, blue, alpha uint32) discord.Color {
	return RGBToColour(RGBAToRGB(red, green, blue, alpha))
}

// RGBToColour takes 3 values and converts it to a Discord colour;
// - red   - 0-255
// - green - 0-255
// - blue  - 0-255
func RGBToColour(red, green, blue uint8) discord.Color {
	var c uint32
	c |= (uint32(red) & 0xff) << 16
	c |= (uint32(green) & 0xff) << 8
	c |= (uint32(blue) & 0xff)

	return discord.Color(c)
}
