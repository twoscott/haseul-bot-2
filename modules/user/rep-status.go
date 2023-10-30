package user

import (
	"fmt"
	"log"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
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
		repsString := util.PluraliseWithCount("rep", remaining)
		message = fmt.Sprintf("You have %s remaining to give.", repsString)
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
