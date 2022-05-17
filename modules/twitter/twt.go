package twitter

import (
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var twtCommand = &router.Command{
	Name:        "twitter",
	Aliases:     []string{"twit", "twt"},
	UseTyping:   true,
	Run:         twtRun,
	SubCommands: make(router.CommandMap),
}

func twtRun(ctx router.CommandCtx, _ []string) {
	dctools.TextReplyNoPing(ctx.State, ctx.Msg,
		"For help with Twitter commands, follow this link: "+
			"https://haseulbot.xyz/#twitter",
	)
}
