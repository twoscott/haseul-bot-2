package vlive

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
)

var vliveFeedsClearCommand = &router.SubCommand{
	Name:        "clear",
	Description: "Clears all VLIVE feeds from the server",
	Handler: &router.CommandHandler{
		Executor: vliveFeedClearExec,
	},
}

func vliveFeedClearExec(ctx router.CommandCtx) {
	_, err := db.VLIVE.ClearGuildMentions(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while clearing all VLIVE roles from " +
				"the database.",
		)
		return
	}

	cleared, err := db.VLIVE.ClearGuildFeeds(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while clearing all VLIVE feeds from " +
				"the database.",
		)
		return
	}
	if cleared == 0 {
		ctx.RespondWarning(
			"There are no VLIVE feeds to be cleared from this server.",
		)
		return
	}

	ctx.RespondSuccess("VLIVE feeds have been cleared from this server.")
}
