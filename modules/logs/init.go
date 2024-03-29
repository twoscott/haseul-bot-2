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
	logsMemberCommand.AddSubCommand(logsMemberChannelCommand)
	logsMemberCommand.AddSubCommand(logsMemberDisableCommand)

	logsCommand.AddSubCommandGroup(logsMessageCommand)
	logsMessageCommand.AddSubCommand(logsMessageChannelCommand)
	logsMessageCommand.AddSubCommand(logsMessageDisableCommand)

	logsCommand.AddSubCommandGroup(logsWelcomeCommand)
	logsWelcomeCommand.AddSubCommand(logsWelcomeChannelCommand)
	logsWelcomeCommand.AddSubCommand(logsWelcomeColourCommand)
	logsWelcomeCommand.AddSubCommand(logsWelcomeDisableCommand)
	logsWelcomeCommand.AddSubCommand(logsWelcomeMessageCommand)
	logsWelcomeCommand.AddSubCommand(logsWelcomeTitleCommand)

	rt.AddMemberJoinHandler(logNewMember)
	rt.AddMemberLeaveHandler(logMemberLeave)
	rt.AddMessageDeleteHandler(logMessageDelete)
	rt.AddMessageUpdateHandler(logMessageUpdate)
}
