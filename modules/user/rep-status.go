package user

import (
	"fmt"
	"log"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var repStatusCommand = &router.SubCommand{
	Name: "status",
	Description: "Returns how many reps you have remaining, or when you will " +
		"be able to give a rep again",
	Handler: &router.CommandHandler{
		Executor: repStatusExec,
		Defer:    true,
	},
}

func repStatusExec(ctx router.CommandCtx) {
	remaining, err := db.Reps.GetUserRepsRemaining(ctx.Interaction.SenderID())
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching remaining reps.")
		return
	}

	var message string
	if remaining > 0 {
		message = fmt.Sprintf("You have %d reps remaining to give.", remaining)
	} else {
		resetTime := getNextRepResetFromNow()
		timeString := dctools.UnixTimestampStyled(
			resetTime,
			dctools.RelativeTime,
		)
		message = fmt.Sprintf(
			"You have no reps remaining to give. "+
				"Your reps will be replenished %s",
			timeString,
		)
	}

	ctx.RespondText(message)
}
