package notifications

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var notiClearCommand = &router.Command{
	Name:      "clear",
	UseTyping: true,
	Run:       notiClearRun,
}

func notiClearRun(ctx router.CommandCtx, _ []string) {
	cleared, err := db.Notifications.Clear(ctx.Msg.Author.ID, ctx.Msg.GuildID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while clearing all notifications from "+
				"the database.",
		)
		return
	}
	if cleared == 0 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"You have no notifications to be cleared in this server.",
		)
		return
	}

	dctools.ReplyWithSuccess(ctx.State, ctx.Msg,
		"Your notifications have been cleared from this server.",
	)
}
