package notification

import (
	"github.com/twoscott/haseul-bot-2/router"
)

const (
	serverScope int64 = iota
	globalScope
)

var notificationCommand = &router.Command{
	Name:        "notification",
	Description: "Commands pertaining to keyword notifications",
}
