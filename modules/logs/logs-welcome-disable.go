package logs

import (
	"fmt"
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
	disabled, err := db.Guilds.DisableWelcomeLogs(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while disabling welcome messages.")
		return
	}

	if !disabled {
		err := fmt.Errorf(
			"welcome logs weren't disabled for %d",
			ctx.Interaction.GuildID,
		)
		log.Println(err)
		ctx.RespondError("Error occurred while disabling welcome messages.")
		return
	}

	ctx.RespondSuccess("Welcome logs disabled.")
}
