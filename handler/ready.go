package handler

import (
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/twoscott/haseul-bot-2/config"
)

var started = false

func (h *Handler) Ready(ev *gateway.ReadyEvent) {
	cfg := config.GetInstance()
	logChannelID := cfg.Discord.LogChannelID
	if logChannelID.IsValid() {
		h.State.SendMessage(logChannelID, "Ready to *Go!~*")
	}

	if !started {
		h.HandleStartupEvent(ev)
		started = true
	}
}
