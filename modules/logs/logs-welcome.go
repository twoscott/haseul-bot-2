package logs

import (
	"github.com/twoscott/haseul-bot-2/router"
)

var logsWelcomeCommand = &router.SubCommandGroup{
	Name:        "welcome",
	Description: "Commands pertaining to member welcome logs",
}
