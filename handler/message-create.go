package handler

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

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
