package handler

import (
	"log"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/twoscott/haseul-bot-2/utils/botutil"
)

var started = false

func (h *Handler) Ready(ev *gateway.ReadyEvent) {
	_, err := botutil.LogText(h.Router.State, "Ready to *Go!~*")
	if err != nil {
		log.Println(err)
	}

	if !started {
		h.Router.HandleStartupEvent(ev)
		started = true
	}
}
