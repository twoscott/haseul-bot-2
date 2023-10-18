package logs

import (
	"fmt"
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
	disabled, err := db.Guilds.DisableMemberLogs(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while disabling member logs.")
		return
	}

	if !disabled {
		err := fmt.Errorf(
			"member logs weren't disabled for %d",
			ctx.Interaction.GuildID,
		)
		log.Println(err)
		ctx.RespondError("Error occurred while disabling member logs.")
		return
	}

	ctx.RespondSuccess("Member logs disabled.")
}
