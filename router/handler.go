// Package handler handles events from Discord's API.
package router

import (
	"log"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/utils/botutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

// Handler wraps router and handles events from the API, and passes them on
// to the router.
type Handler struct {
	db      *database.DB
	Router  *Router
	Started bool
}

const defaultPrefix = "."

// New returns a new instance of Handler.
func NewHandler(router *Router) *Handler {
	return &Handler{
		db:      database.GetInstance(),
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
	case *discord.StringSelectInteraction:
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

	go h.Router.HandleMessage(msg.Message, msg.Member)

	if len(msg.Content) == 0 {
		return
	}
	args := strings.Fields(msg.Content[1:])
	if len(args) < 1 {
		return
	}

	prefix, _ := h.db.Guilds.GetLegacyPrefix(msg.GuildID)
	if prefix == "" {
		prefix = defaultPrefix
	}

	if !strings.HasPrefix(msg.Content, string(prefix)) {
		return
	}

	h.Router.HandleLegacyCommand(msg, args)
}

func (h *Handler) MessageDelete(ev *gateway.MessageDeleteEvent) {
	msg, err := h.Router.State.Message(ev.ChannelID, ev.ID)
	if err != nil {
		log.Println(err)
		return
	}

	go h.Router.HandleMessageDelete(*msg)
}

func (h *Handler) MessageUpdate(ev *gateway.MessageUpdateEvent) {
	old, err := h.Router.State.Message(ev.ChannelID, ev.ID)
	if err != nil {
		log.Println(err)
		return
	}

	go h.Router.HandleMessageUpdate(*old, ev.Message, ev.Member)
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
