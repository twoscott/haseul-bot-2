package notifications

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var notiGlobClearCommand = &router.Command{
	Name:      "clear",
	UseTyping: true,
	Run:       notiGlobClearRun,
}

func notiGlobClearRun(ctx router.CommandCtx, _ []string) {
	cleared, err := db.Notifications.ClearGlobal(ctx.Msg.Author.ID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while clearing all global notifications.",
		)
		return
	}
	if cleared == 0 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"You have no global notifications to be cleared.",
		)
		return
	}

	dctools.ReplyWithSuccess(ctx.State, ctx.Msg,
		"Your global notifications have been cleared.",
	)
}
