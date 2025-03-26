package admin

import "github.com/twoscott/haseul-bot-2/router"

func Init(rt *router.Router) {
	rt.AddCommand(adminCommand)
	adminCommand.AddSubCommandGroup(adminServer)
	adminServer.AddSubCommand(adminServerList)
	adminServer.AddSubCommand(adminServerInfoCommand)
}
