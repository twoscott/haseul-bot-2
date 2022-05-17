package twitter

import (
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var twtToggleCommand = &router.Command{
	Name:        "toggle",
	UseTyping:   true,
	Run:         twtToggleRun,
	SubCommands: make(router.CommandMap),
}

func twtToggleRun(ctx router.CommandCtx, _ []string) {
	dctools.TextReplyNoPing(ctx.State, ctx.Msg,
		"For help with Twitter commands, follow this link: "+
			"https://haseulbot.xyz/#twitter",
	)
}
