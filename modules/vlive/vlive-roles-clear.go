package vlive

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
)

var vliveRolesClearCommand = &router.SubCommand{
	Name:        "clear",
	Description: "Clears all mention roles for all VLIVE feeds",
	Handler: &router.CommandHandler{
		Executor: vliveRoleClearExec,
	},
}

func vliveRoleClearExec(ctx router.CommandCtx) {
	cleared, err := db.VLIVE.ClearGuildMentions(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while clearing all VLIVE roles " +
				"from the database.",
		)
		return
	}
	if cleared == 0 {
		ctx.RespondWarning(
			"There are no VLIVE roles to be cleared in this server.",
		)
		return
	}

	ctx.RespondSuccess("VLIVE roles have been cleared from this server.")
}
