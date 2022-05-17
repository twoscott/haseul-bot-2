package handler

import (
	"log"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/twoscott/haseul-bot-2/cache"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

const defaultPrefix = "."

func (h *Handler) MessageCreate(msg *gateway.MessageCreateEvent) {
	if msg.Author.Bot {
		return
	}
	if !dctools.IsUserMessage(msg.Type) {
		return
	}

	channel, err := h.State.Channel(msg.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}
	if channel.Type == discord.DirectMessage {
		return
	}

	go h.HandleMessage(msg)

	args := strings.Fields(msg.Content[1:])
	if len(args) < 1 {
		return
	}

	prefix := cache.GetInstance().GetPrefix(msg.GuildID)
	if prefix == "" {
		prefix = defaultPrefix
	}

	if !strings.HasPrefix(msg.Content, prefix) {
		return
	}

	h.HandleCommand(msg, args)
}
