package vlive

import (
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var vlvFeedCommand = &router.Command{
	Name:        "feed",
	Aliases:     []string{"feeds"},
	UseTyping:   true,
	Run:         vlvFeedRun,
	SubCommands: make(router.CommandMap),
}

func vlvFeedRun(ctx router.CommandCtx, _ []string) {
	dctools.TextReplyNoPing(ctx.State, ctx.Msg,
		"For help with VLIVE commands, follow this link: "+
			"https://haseulbot.xyz/#vlive",
	)
}
