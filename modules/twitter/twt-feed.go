package twitter

import (
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var twtFeedCommand = &router.Command{
	Name:        "feed",
	Aliases:     []string{"feeds"},
	UseTyping:   true,
	Run:         twtFeedRun,
	SubCommands: make(router.CommandMap),
}

func twtFeedRun(ctx router.CommandCtx, _ []string) {
	dctools.TextReplyNoPing(ctx.State, ctx.Msg,
		"For help with Twitter commands, follow this link: "+
			"https://haseulbot.xyz/#twitter",
	)
}
