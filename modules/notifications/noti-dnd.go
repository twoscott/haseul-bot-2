package notifications

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var notiDndCommand = &router.Command{
	Name:      "dnd",
	Aliases:   []string{"donotdisturb", "sleep"},
	UseTyping: true,
	Run:       notiDndRun,
}

func notiDndRun(ctx router.CommandCtx, _ []string) {
	dndOn, err := db.Notifications.ToggleDnD(ctx.Msg.Author.ID)
	if err != nil {
		log.Println(err)
		dctools.SendError(ctx.State, ctx.Msg.ChannelID,
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

	dctools.SendSuccess(ctx.State, ctx.Msg.ChannelID,
		"Your do not disturb status was turned "+status+".",
	)
}
