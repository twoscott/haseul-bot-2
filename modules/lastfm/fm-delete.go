package lastfm

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var fmDeleteCommand = &router.Command{
	Name:      "delete",
	Aliases:   []string{"remove"},
	UseTyping: true,
	Run:       fmDeleteRun,
}

func fmDeleteRun(ctx router.CommandCtx, _ []string) {
	del, err := db.LastFM.DeleteUser(ctx.Msg.Author.ID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while trying to delete your Last.fm username",
		)
		return
	}
	if !del {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"You don't have a Last.fm username linked to your Discord account.",
		)
		return
	}

	dctools.ReplyWithSuccess(ctx.State, ctx.Msg,
		"Your Last.fm username was deleted from the records.",
	)
}
