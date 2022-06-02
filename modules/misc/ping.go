package misc

import (
	"fmt"
	"log"
	"time"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var pingCommand = &router.Command{
	Name:        "ping",
	Description: "Times how long it takes for the bot to respond",
	Handler: &router.CommandHandler{
		Executor: pingExec,
	},
}

func pingExec(ctx router.CommandCtx) {
	start := time.Now()

	err := dctools.DeferResponse(ctx.State, ctx.Interaction)

	ping := time.Since(start)

	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while timing response.")
		return
	}

	dctools.FollowupRespondText(ctx.State, ctx.Interaction,
		fmt.Sprintf("Ping: %dms", ping.Milliseconds()),
	)
}
