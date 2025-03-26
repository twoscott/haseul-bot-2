package admin

import (
	"github.com/twoscott/haseul-bot-2/router"
)

var adminServer = &router.SubCommandGroup{
	Name:        "server",
	Description: "Admin commands for servers",
}
