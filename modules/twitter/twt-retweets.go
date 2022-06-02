package twitter

import (
	"github.com/twoscott/haseul-bot-2/router"
)

var twtRetweetsCommand = &router.SubCommandGroup{
	Name:        "retweets",
	Description: "Commands pertaining to Twitter feed Retweets posting",
}
