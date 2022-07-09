package twitter

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
)

var twtFeedsClearCommand = &router.SubCommand{
	Name:        "clear",
	Description: "Clears all Twitter feeds from the server",
	Handler: &router.CommandHandler{
		Executor: twtFeedClearExec,
	},
}

func twtFeedClearExec(ctx router.CommandCtx) {
	_, err := db.Twitter.ClearGuildMentions(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while clearing all Twitter roles from " +
				"the database.",
		)
		return
	}

	cleared, err := db.Twitter.ClearGuildFeeds(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while clearing all Twitter feeds from " +
				"the database.",
		)
		return
	}
	if cleared == 0 {
		ctx.RespondWarning(
			"There are no Twitter feeds to be cleared from this server.",
		)
		return
	}

	ctx.RespondSuccess("Twitter feeds have been cleared from this server.")
}
