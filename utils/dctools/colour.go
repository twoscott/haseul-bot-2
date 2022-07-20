package dctools

import "github.com/diamondburned/arikawa/v3/discord"

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
