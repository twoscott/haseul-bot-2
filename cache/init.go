package cache

import "github.com/twoscott/haseul-bot-2/router"

func Init(rt *router.Router) {
	c := GetInstance()
	rt.AddStartupListener(c.onStartup)
	rt.AddGuildJoinHandler(c.onGuildJoin)
}
