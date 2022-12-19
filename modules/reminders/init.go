package reminders

import (
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
)

var db *database.DB

func Init(rt *router.Router) {
	db = database.GetInstance()

	rt.AddStartupListener(onStartup)

	rt.AddCommand(remindersCommand)
	remindersCommand.AddSubCommand(remindersAddCommand)
	remindersCommand.AddSubCommand(remindersClearCommand)
	remindersCommand.AddSubCommand(remindersDeleteCommand)
	remindersCommand.AddSubCommand(remindersListCommand)
}

func onStartup(rt *router.Router, _ *gateway.ReadyEvent) {
	startCheckingReminders(rt.State)
}
