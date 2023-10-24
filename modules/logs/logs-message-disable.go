package logs

import (
	"fmt"
	"log"

	"github.com/twoscott/haseul-bot-2/router"
)

var logsMessageDisableCommand = &router.SubCommand{
	Name:        "disable",
	Description: "Stops message logs from being posted to the server",
	Handler: &router.CommandHandler{
		Executor: logsMessageDisableExec,
	},
}

func logsMessageDisableExec(ctx router.CommandCtx) {
	disabled, err := db.Guilds.DisableMessageLogs(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while disabling message logs.")
		return
	}

	if !disabled {
		err := fmt.Errorf(
			"message logs weren't disabled for %d",
			ctx.Interaction.GuildID,
		)
		log.Println(err)
		ctx.RespondError("Error occurred while disabling message logs.")
		return
	}

	ctx.RespondSuccess("Message logs disabled.")
}
