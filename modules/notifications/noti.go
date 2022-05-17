package notifications

import (
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var notiCommand = &router.Command{
	Name:        "notification",
	Aliases:     []string{"notifications", "notif", "noti"},
	UseTyping:   true,
	Run:         notiRun,
	SubCommands: make(router.CommandMap),
}

func notiRun(ctx router.CommandCtx, _ []string) {
	dctools.TextReplyNoPing(ctx.State, ctx.Msg,
		"For help with Notification commands, follow this link: "+
			"https://haseulbot.xyz/#notifications",
	)
}
