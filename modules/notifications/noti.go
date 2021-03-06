package notifications

import (
	"github.com/twoscott/haseul-bot-2/router"
)

const (
	serverScope int64 = iota
	globalScope
)

var notiCommand = &router.Command{
	Name:        "notifications",
	Description: "Commands pertaining to keyword notifications",
}
