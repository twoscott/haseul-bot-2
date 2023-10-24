package logs

import (
	"github.com/twoscott/haseul-bot-2/router"
)

var logsMessageCommand = &router.SubCommandGroup{
	Name:        "message",
	Description: "Commands pertaining to message logs",
}
