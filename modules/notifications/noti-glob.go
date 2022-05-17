package notifications

import (
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var notiGlobCommand = &router.Command{
	Name:        "global",
	Aliases:     []string{"glob"},
	UseTyping:   true,
	Run:         notiGlobRun,
	SubCommands: make(router.CommandMap),
}

func notiGlobRun(ctx router.CommandCtx, _ []string) {
	dctools.TextReplyNoPing(ctx.State, ctx.Msg,
		"For help with Notification commands, follow this link: "+
			"https://haseulbot.xyz/#notifications",
	)
}
