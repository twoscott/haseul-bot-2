package twitter

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
)

var twtRolesClearCommand = &router.SubCommand{
	Name:        "clear",
	Description: "Clears all mention roles for all Twitter feeds",
	Handler: &router.CommandHandler{
		Executor: twtRoleClearExec,
	},
}

func twtRoleClearExec(ctx router.CommandCtx) {
	cleared, err := db.Twitter.ClearGuildMentions(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while clearing all Twitter roles " +
				"from the database.",
		)
		return
	}
	if cleared == 0 {
		ctx.RespondWarning(
			"There are no Twitter roles to be cleared in this server.",
		)
		return
	}

	ctx.RespondSuccess("Twitter roles have been cleared from this server.")
}
