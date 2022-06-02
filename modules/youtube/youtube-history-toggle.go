package youtube

import (
	"github.com/twoscott/haseul-bot-2/router"
)

var youTubeHistoryToggleCommand = &router.SubCommand{
	Name:        "toggle",
	Description: "Toggles the tracking ofs YouTube search history",
	Handler: &router.CommandHandler{
		Executor: youTubeHistoryToggleExec,
	},
}

func youTubeHistoryToggleExec(ctx router.CommandCtx) {
	err := db.YouTube.ToggleHistory(ctx.Interaction.SenderID())
	if err != nil {
		ctx.RespondError("Error occurred toggling search history.")
		return
	}

	toggle, err := db.YouTube.GetHistoryToggle(ctx.Interaction.SenderID())
	if err != nil {
		ctx.RespondSuccess("Search history tracking was toggled.")
		return
	}

	switch toggle {
	case true:
		ctx.RespondSuccess(
			"Your YouTube search history will now be tracked.",
		)
	case false:
		db.YouTube.ClearHistory(ctx.Interaction.SenderID())
		ctx.RespondSuccess(
			"Your YouTube search history will no longer be tracked.",
		)
	}
}
