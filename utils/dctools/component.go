package dctools

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

// DisabledButtons returns a copy of the provided buttons, all disabled.
func DisabledButtons(
	buttons discord.ActionRowComponent) discord.ActionRowComponent {

	newButtons := make(discord.ActionRowComponent, len(buttons))
	copy(newButtons, buttons)

	for i, c := range newButtons {
		switch b := c.(type) {
		case *discord.ButtonComponent:
			newButton := *b
			newButton.Disabled = true
			newButtons[i] = &newButton
		}
	}

	return newButtons
}
