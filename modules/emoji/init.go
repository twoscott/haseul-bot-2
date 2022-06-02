package emoji

import "github.com/twoscott/haseul-bot-2/router"

func Init(rt *router.Router) {
	rt.AddCommand(emojiCommand)
	emojiCommand.AddSubCommand(emojiExpandCommand)

	// TODO:
	//	/emoji list
	// 	/emoji upload [name] [attachment or url] - automatically resize the image
	// 	/emoji delete [emoji]
}
