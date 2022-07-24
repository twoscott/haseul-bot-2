package logs

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
)

var logsWelcomeDisableCommand = &router.SubCommand{
	Name:        "disable",
	Description: "Stops member logs from being posted to the server",
	Handler: &router.CommandHandler{
		Executor: logsWelcomeDisableExec,
	},
}

func logsWelcomeDisableExec(ctx router.CommandCtx) {
	_, err := db.Guilds.DisableWelcomeLogs(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while disabling member logs.")
		return
	}

	ctx.RespondSuccess("Welcome logs disabled.")
}
