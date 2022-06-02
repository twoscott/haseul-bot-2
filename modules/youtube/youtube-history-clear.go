package youtube

import (
	"fmt"

	"github.com/twoscott/haseul-bot-2/router"
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
		fmt.Sprintf("Cleared %d entries from YouTube search history.", cleared),
	)
}
