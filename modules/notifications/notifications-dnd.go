package notifications

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
)

var notificationsDndCommand = &router.SubCommand{
	Name: "do-not-disturb",
	Description: "Toggles whether Do Not Disturb is turned on " +
		"for notifications",
	Handler: &router.CommandHandler{
		Executor:  notificationsDndExec,
		Ephemeral: true,
	},
}

func notificationsDndExec(ctx router.CommandCtx) {
	dndOn, err := db.Notifications.ToggleDnD(ctx.Interaction.SenderID())
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while toggling your do not disturb status.",
		)
		return
	}

	var status string
	if dndOn {
		status = "on"
	} else {
		status = "off"
	}

	ctx.RespondSuccess(
		"Your do not disturb status was turned " + status + ".",
	)
}
