package information

import "github.com/twoscott/haseul-bot-2/router"

func Init(rt *router.Router) {
	rt.MustRegisterCommand(avatarCommand)
	rt.MustRegisterCommand(bannerCommand)
	rt.MustRegisterCommand(guildbannerCommand)
	rt.MustRegisterCommand(botCommand)
	rt.MustRegisterCommand(cacheCommand)
	rt.MustRegisterCommand(guildCommand)
	rt.MustRegisterCommand(iconCommand)
	rt.MustRegisterCommand(userCommand)
}
