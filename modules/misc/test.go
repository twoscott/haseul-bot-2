package misc

import (
	"github.com/twoscott/haseul-bot-2/router"
)

var testCommand = &router.Command{
	Name:      "test",
	UseTyping: true,
	Run:       testRun,
}

func testRun(ctx router.CommandCtx, args []string) {}
