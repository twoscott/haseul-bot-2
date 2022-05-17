package misc

import (
	"fmt"
	"log"
	"time"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var pingCommand = &router.Command{
	Name: "ping",
	Run:  pingRun,
}

func pingRun(ctx router.CommandCtx, args []string) {
	start := time.Now()

	msg, err := dctools.TextReplyNoPing(ctx.State, ctx.Msg, "Ping: ...ms")
	if err != nil {
		log.Println(err)
		return
	}

	ping := time.Since(start)

	ctx.State.EditText(msg.ChannelID, msg.ID,
		fmt.Sprintf("Ping: %dms", ping.Milliseconds()),
	)
}
