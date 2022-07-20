package dctools

import (
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

const (
	// EmbedFooterSep is a seperator that mimics Discord's separator between
	// footers and timestamps.
	EmbedFooterSep = util.ThreePerEmSpace + "•" + util.ThreePerEmSpace
	// EmbedTimeFormat is a format for representing times in embeds.
	EmbedTimeFormat = "02 Jan 2006 15:04:05"
)

// EmbedTime returns t converted to UTC and string formatted for embeds.
func EmbedTime(t time.Time) string {
	return t.UTC().Format(EmbedTimeFormat)
}

// HexToColour converts a string representing a hexadecimal colour value
// to a Discord integer colour.
// e.g: "#3251cf" > 0x3251cf
//
// String must either begin with a '#' or consist only of hexadecimal values.
func HexToColour(hex string) (discord.Color, error) {
	hex = strings.TrimPrefix(hex, "#")

	colourInt64, err := strconv.ParseInt(hex, 16, 32)
	if err != nil {
		return discord.NullColor, err
	}

	return discord.Color(colourInt64), nil
}

// EmptyEmbedField returns an embed field filled with zero-width spaces.
func EmptyEmbedField() discord.EmbedField {
	return discord.EmbedField{
		Name:   "\u200b",
		Value:  "\u200b",
		Inline: true,
	}
}

// SeparateFooter takes multiple strings as sections to an embed footer, and
// joins them, separated by an embed footer separator.
func SeparateEmbedFooter(sections ...string) string {
	return strings.Join(sections, EmbedFooterSep)
}
