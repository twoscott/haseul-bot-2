package logs

import (
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/inviteutil"
)

var (
	db            *database.DB
	inviteTracker *inviteutil.Tracker
)

func Init(rt *router.Router) {
	db = database.GetInstance()
	inviteTracker = inviteutil.GetTracker()

	rt.AddCommand(logsCommand)
	logsCommand.AddSubCommandGroup(logsMemberCommand)
	logsMemberCommand.AddSubCommand(logsMemberSetCommand)
	logsMemberCommand.AddSubCommand(logsMemberDisableCommand)
	logsCommand.AddSubCommandGroup(logsWelcomeCommand)
	logsWelcomeCommand.AddSubCommand(logsWelcomeSetCommand)
	logsWelcomeCommand.AddSubCommand(logsWelcomeDisableCommand)

	rt.AddMemberJoinHandler(logNewMember)
	rt.AddMemberLeaveHandler(logMemberLeave)
}
