package notifications

import (
	"github.com/twoscott/haseul-bot-2/router"
)

const (
	serverScope = iota
	globalScope
)

var notificationsCommand = &router.Command{
	Name:        "notifications",
	Description: "Commands pertaining to keyword notifications",
}
