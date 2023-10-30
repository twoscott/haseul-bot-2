package dctools

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
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
	EmbedBackColour = 0x2b2d31
)

var ErrInvalidHexColour = errors.New("invalid hexadecimal RGB value provided")

// ColourInvalid returns whether a colour is either null or a default colour.
func ColourInvalid(colour discord.Color) bool {
	return colour == 0x000000 ||
		colour == discord.DefaultEmbedColor ||
		colour == discord.NullColor
}

// EmbedColour returns a Discord colour where default black colours are
// replaced with the default embed background colour.
func EmbedColour(colour discord.Color) discord.Color {
	if ColourInvalid(colour) {
		colour = EmbedBackColour
	}
	return colour
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
