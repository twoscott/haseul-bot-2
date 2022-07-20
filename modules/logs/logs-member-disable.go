package logs

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
)

var logsMemberDisableCommand = &router.SubCommand{
	Name:        "disable",
	Description: "Stops member logs from being posted to the server",
	Handler: &router.CommandHandler{
		Executor: logsMemberDisableExec,
	},
}

func logsMemberDisableExec(ctx router.CommandCtx) {
	_, err := db.Guilds.DisableMemberLogs(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while disabling member logs.")
		return
	}

	ctx.RespondSuccess("Member logs disabled.")
}
