package dctools

import (
	"errors"
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
	EmbedBackColour = 0x2f3136
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
