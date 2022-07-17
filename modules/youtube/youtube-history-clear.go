package youtube

import (
	"fmt"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var youTubeHistoryClearCommand = &router.SubCommand{
	Name:        "clear",
	Description: "Clears YouTube search history",
	Handler: &router.CommandHandler{
		Executor: youTubeHistoryClearExec,
	},
}

func youTubeHistoryClearExec(ctx router.CommandCtx) {
	cleared, err := db.YouTube.ClearHistory(ctx.Interaction.SenderID())
	if err != nil {
		ctx.RespondError(
			"Error occurred trying to clear search history.",
		)
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf(
			"Cleared %d %s from YouTube search history.",
			cleared,
			util.Pluralise("entry", cleared),
		),
	)
}
