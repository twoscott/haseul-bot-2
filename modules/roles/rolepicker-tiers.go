package roles

import (
	"github.com/twoscott/haseul-bot-2/router"
)

var rolePickerTiers = &router.SubCommandGroup{
	Name:        "tiers",
	Description: "Commands pertaining to role picker roles",
}
