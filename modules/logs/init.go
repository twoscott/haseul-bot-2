package logs

import (
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
)

var db *database.DB

func Init(rt *router.Router) {
	db = database.GetInstance()

	rt.AddCommand(logsCommand)
	logsCommand.AddSubCommandGroup(logsMemberCommand)
	logsMemberCommand.AddSubCommand(logsMemberSetCommand)
	logsMemberCommand.AddSubCommand(logsMemberDisableCommand)

	rt.AddMemberJoinHandler(logMemberJoin)
	rt.AddMemberLeaveHandler(logMemberLeave)
}
