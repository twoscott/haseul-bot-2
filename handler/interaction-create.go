package handler

import (
	"github.com/diamondburned/arikawa/v3/gateway"
)

func (h *Handler) InteractionCreate(
	interaction *gateway.InteractionCreateEvent) {

	switch interaction.Type {
	case gateway.ButtonInteraction:
		h.HandleButton(interaction)
	}
}
