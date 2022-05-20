package dctools

import (
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
	// EmbedBackColour is the default embed background colour Discord uses.
	EmbedBackColour = 0x2f3136
)

// EmbedTime returns t converted to UTC and string formatted for embeds.
func EmbedTime(t time.Time) string {
	return t.UTC().Format(EmbedTimeFormat)
}

// EmbedColour returns a Discord colour where default black colours are
// replaced with the default embed background colour.
func EmbedColour(colour discord.Color) discord.Color {
	if colour == 0x000000 ||
		colour == discord.DefaultEmbedColor ||
		colour == discord.NullColor {
		colour = EmbedBackColour
	}
	return colour
}

// EmptyEmbedField returns an embed field filled with zero-width spaces.
func EmptyEmbedField() discord.EmbedField {
	return discord.EmbedField{
		Name:   "\u200b",
		Value:  "\u200b",
		Inline: true,
	}
}
