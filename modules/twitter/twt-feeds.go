package twitter

import (
	"github.com/twoscott/haseul-bot-2/router"
)

var twtFeedsCommand = &router.SubCommandGroup{
	Name:        "feeds",
	Description: "Commands pertaining to Twitter feeds",
}
