// Package handler handles events from Discord's API.
package router

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/twoscott/haseul-bot-2/utils/botutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

// Handler wraps router and handles events from the API, and passes them on
// to the router.
type Handler struct {
	Router  *Router
	Started bool
}

// New returns a new instance of Handler.
func NewHandler(router *Router) *Handler {
	return &Handler{
		Router:  router,
		Started: false,
	}
}

func (h *Handler) GuildJoin(guild *state.GuildJoinEvent) {
	h.Router.HandleGuildJoin(guild)
}

func (h *Handler) InteractionCreate(
	interaction *gateway.InteractionCreateEvent) {

	switch data := interaction.Data.(type) {
	case *discord.CommandInteraction:
		h.Router.HandleCommand(&interaction.InteractionEvent, data)
	case *discord.AutocompleteInteraction:
		h.Router.HandleAutocomplete(&interaction.InteractionEvent, data)
	case *discord.ButtonInteraction:
		h.Router.HandleButtonPress(&interaction.InteractionEvent, data)
	case *discord.ModalInteraction:
		h.Router.HandleModalSubmit(&interaction.InteractionEvent, data)
	case *discord.SelectInteraction:
		h.Router.HandleSelect(&interaction.InteractionEvent, data)
	}
}

func (h *Handler) MemberJoin(ev *gateway.GuildMemberAddEvent) {
	h.Router.HandleMemberJoin(ev.Member, ev.GuildID)
}

func (h *Handler) MemberLeave(ev *gateway.GuildMemberRemoveEvent) {
	h.Router.HandleMemberLeave(ev.User, ev.GuildID)
}

func (h *Handler) MessageCreate(msg *gateway.MessageCreateEvent) {
	if msg.Author.Bot {
		return
	}
	if !dctools.IsUserMessage(msg.Type) {
		return
	}

	channel, err := h.Router.State.Channel(msg.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}
	if channel.Type == discord.DirectMessage {
		return
	}

	h.Router.HandleMessage(msg.Message, msg.Member)
}

func (h *Handler) MessageDelete(ev *gateway.MessageDeleteEvent) {
	msg, err := h.Router.State.Message(ev.ChannelID, ev.ID)
	if err != nil {
		log.Println(err)
		return
	}

	h.Router.HandleMessageDelete(*msg)
}

func (h *Handler) MessageUpdate(ev *gateway.MessageUpdateEvent) {
	old, err := h.Router.State.Message(ev.ChannelID, ev.ID)
	if err != nil {
		log.Println(err)
		return
	}

	h.Router.HandleMessageUpdate(*old, ev.Message, ev.Member)
}

func (h *Handler) Ready(ev *gateway.ReadyEvent) {
	_, err := botutil.LogText(h.Router.State, "Ready to *Go!~*")
	if err != nil {
		log.Println(err)
	}

	if !h.Started {
		h.Router.HandleStartupEvent(ev)
		h.Started = true
	}
}
