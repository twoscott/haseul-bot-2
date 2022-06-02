package server

import "github.com/twoscott/haseul-bot-2/router"

func Init(rt *router.Router) {
	rt.AddCommand(serverCommand)
	serverCommand.AddSubCommand(serverBannerCommand)
	serverCommand.AddSubCommand(serverIconCommand)
	serverCommand.AddSubCommand(serverInfoCommand)
}
