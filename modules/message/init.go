package message

import (
	"github.com/twoscott/haseul-bot-2/router"
)

func Init(rt *router.Router) {
	rt.AddCommand(messageCommand)
	messageCommand.AddSubCommand(messageSendCommand)
	messageCommand.AddSubCommand(messageEditCommand)
	messageCommand.AddSubCommand(messageFetchCommand)

	rt.AddCommand(editMessageCommand)
}
