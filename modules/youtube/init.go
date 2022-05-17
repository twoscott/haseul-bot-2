package youtube

import "github.com/twoscott/haseul-bot-2/router"

func Init(router *router.Router) {
	router.MustRegisterCommand(ytCommand)
}
