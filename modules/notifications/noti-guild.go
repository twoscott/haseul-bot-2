package notifications

import (
	"github.com/twoscott/haseul-bot-2/router"
)

var notiGuildCommand = &router.Command{
	Name:        "guild",
	Aliases:     []string{"server"},
	SubCommands: make(router.CommandMap),
}
