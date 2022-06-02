package handler

import (
	"log"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/twoscott/haseul-bot-2/config"
)

var started = false

func (h *Handler) Ready(ev *gateway.ReadyEvent) {
	cfg := config.GetInstance()
	logChannelID := cfg.Discord.LogChannelID

	var err error
	if !logChannelID.IsValid() {
		log.Println("Configured log channel ID is invalid.")
	} else {
		_, err = h.Router.State.SendMessage(logChannelID, "Ready to *Go!~*")
	}

	if err != nil {
		log.Println(err)
	}

	if !started {
		h.Router.HandleStartupEvent(ev)
		started = true
	}
}
