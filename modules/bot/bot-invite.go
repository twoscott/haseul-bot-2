package bot

import (
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/botutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var botInviteCommand = &router.SubCommand{
	Name:        "invite",
	Description: "Sends a link to invite Haseul Bot to your server",
	Handler: &router.CommandHandler{
		Executor: botInviteExec,
	},
}

func botInviteExec(ctx router.CommandCtx) {
	msg := dctools.Hyperlink(
		"Invite Haseul Bot to your server",
		botutil.Invite,
	)
	ctx.RespondText(dctools.Bold(msg))
}
