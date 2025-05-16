package lastfm

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
)

var lastFMDeleteCommand = &router.SubCommand{
	Name:        "delete",
	Description: "Deletes your Last.fm username from the records",
	Handler: &router.CommandHandler{
		Executor: lastFMDeleteExec,
	},
}

func lastFMDeleteExec(ctx router.CommandCtx) {
	del, err := db.LastFM.DeleteUser(ctx.Interaction.SenderID())
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while trying to delete your Last.fm username",
		)
		return
	}
	if !del {
		ctx.RespondWarning(
			"You don't have a Last.fm username linked to your Discord account.",
		)
		return
	}

	ctx.RespondSuccess(
		"Your Last.fm username was deleted from the records.",
	)
}
