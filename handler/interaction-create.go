package handler

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func (h *Handler) InteractionCreate(
	interaction *gateway.InteractionCreateEvent) {

	switch data := interaction.Data.(type) {
	case *discord.ButtonInteraction:
		h.HandleButton(interaction, data)
	}
}
