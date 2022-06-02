package user

import "github.com/twoscott/haseul-bot-2/router"

func Init(rt *router.Router) {
	rt.AddCommand(userCommand)
	userCommand.AddSubCommand(userAvatarCommand)
	userCommand.AddSubCommand(userBannerCommand)
	userCommand.AddSubCommand(userInfoCommand)
}
