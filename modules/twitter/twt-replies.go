package twitter

import (
	"github.com/twoscott/haseul-bot-2/router"
)

var twtRepliesCommand = &router.SubCommandGroup{
	Name:        "replies",
	Description: "Commands pertaining to Twitter feed replies posting",
}
