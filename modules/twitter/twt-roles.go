package twitter

import (
	"github.com/twoscott/haseul-bot-2/router"
)

var twtRolesCommand = &router.SubCommandGroup{
	Name:        "roles",
	Description: "Commands pertaining to Twitter feed mention roles",
}
