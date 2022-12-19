package reminders

import (
	"fmt"
	"log"

	"github.com/twoscott/haseul-bot-2/router"
)

var remindersClearCommand = &router.SubCommand{
	Name:        "clear",
	Description: "Delete all reminders you have set",
	Handler: &router.CommandHandler{
		Executor:  remindersClearExec,
		Ephemeral: true,
	},
}

func remindersClearExec(ctx router.CommandCtx) {
	cleared, err := db.Reminders.ClearByUser(ctx.Interaction.SenderID())
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while trying to delete reminder.")
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf("Deleted %d reminders.", cleared),
	)
}
