package handler

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func (h *Handler) InteractionCreate(
	interaction *gateway.InteractionCreateEvent) {

	switch data := interaction.Data.(type) {
	case *discord.AutocompleteInteraction:
		h.Router.HandleAutocomplete(&interaction.InteractionEvent, data)
	case *discord.ButtonInteraction:
		h.Router.HandleButtonPress(&interaction.InteractionEvent, data)
	case *discord.CommandInteraction:
		h.Router.HandleCommand(&interaction.InteractionEvent, data)
	}
}
