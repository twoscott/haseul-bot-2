package misc

import (
	"github.com/twoscott/haseul-bot-2/router"
)

var testCommand = &router.Command{
	Name:        "test",
	Description: "Test command",
	Handler: &router.CommandHandler{
		Executor: testExec,
	},
}

func testExec(ctx router.CommandCtx) {}
