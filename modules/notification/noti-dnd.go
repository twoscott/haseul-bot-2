package notification

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
)

var notificationDndCommand = &router.SubCommand{
	Name: "dnd",
	Description: "Toggles whether Do Not Disturb is turned on " +
		"for notifications",
	Handler: &router.CommandHandler{
		Executor:  notificationDndExec,
		Ephemeral: true,
	},
}

func notificationDndExec(ctx router.CommandCtx) {
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
