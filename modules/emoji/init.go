package emoji

import "github.com/twoscott/haseul-bot-2/router"

func Init(rt *router.Router) {
	rt.MustRegisterCommand(emojiCommand)
}
