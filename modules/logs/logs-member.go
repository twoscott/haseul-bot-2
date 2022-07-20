package logs

import (
	"github.com/twoscott/haseul-bot-2/router"
)

var logsMemberCommand = &router.SubCommandGroup{
	Name:        "member",
	Description: "Commands pertaining to member logs",
}
