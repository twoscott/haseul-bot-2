package twitter

import (
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var twtRoleCommand = &router.Command{
	Name:        "role",
	Aliases:     []string{"roles"},
	UseTyping:   true,
	Run:         twtRoleRun,
	SubCommands: make(router.CommandMap),
}

func twtRoleRun(ctx router.CommandCtx, _ []string) {
	dctools.TextReplyNoPing(ctx.State, ctx.Msg,
		"For help with Twitter commands, follow this link: "+
			"https://haseulbot.xyz/#twitter",
	)
}
