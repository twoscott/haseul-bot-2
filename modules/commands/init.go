package commands

import (
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
)

var db *database.DB

func Init(rt *router.Router) {
	db = database.GetInstance()

	rt.AddCommand(commandsCommand)
	commandsCommand.AddSubCommand(commandsAddCommand)
	commandsCommand.AddSubCommand(commandsInfoCommand)
	commandsCommand.AddSubCommand(commandsListCommand)
	commandsCommand.AddSubCommand(commandsDeleteCommand)
	commandsCommand.AddSubCommand(commandsSearchCommand)
	commandsCommand.AddSubCommand(commandsUseCommand)
}
