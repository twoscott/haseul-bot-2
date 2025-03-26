package admin

import (
	"github.com/twoscott/haseul-bot-2/router"
)

var adminCommand = &router.Command{
	Name:        "admin",
	Description: "Admin commands",
	IsAdminOnly: true,
}
