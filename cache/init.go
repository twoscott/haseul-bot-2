package cache

import "github.com/twoscott/haseul-bot-2/router"

func Init(rt *router.Router) {
	c := GetInstance()
	rt.RegisterStartupListener(c.onStartup)
	rt.RegisterGuildJoinHandler(c.onGuildJoin)
}
