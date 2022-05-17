package vlive

import (
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var vlvCommand = &router.Command{
	Name:        "vlive",
	Aliases:     []string{"vlv"},
	UseTyping:   true,
	Run:         vlvRun,
	SubCommands: make(router.CommandMap),
}

func vlvRun(ctx router.CommandCtx, _ []string) {
	dctools.TextReplyNoPing(ctx.State, ctx.Msg,
		"For help with VLIVE commands, follow this link: "+
			"https://haseulbot.xyz/#vlive",
	)
}
