package notifications

import (
	"github.com/twoscott/haseul-bot-2/router"
)

var notiChannelCommand = &router.Command{
	Name:        "channel",
	SubCommands: make(router.CommandMap),
}
