package command

import (
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
)

var db *database.DB

func Init(rt *router.Router) {
	db = database.GetInstance()

	rt.AddCommand(commandCommand)
	commandCommand.AddSubCommand(commandAddCommand)
	commandCommand.AddSubCommand(commandInfoCommand)
	commandCommand.AddSubCommand(commandListCommand)
	commandCommand.AddSubCommand(commandDeleteCommand)
	commandCommand.AddSubCommand(commandSearchCommand)
	commandCommand.AddSubCommand(commandUseCommand)
}
