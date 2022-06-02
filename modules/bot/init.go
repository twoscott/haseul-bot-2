package bot

import "github.com/twoscott/haseul-bot-2/router"

func Init(rt *router.Router) {
	rt.AddCommand(botCommand)
	botCommand.AddSubCommand(botCacheCommand)
	botCommand.AddSubCommand(botInfoCommand)
}
