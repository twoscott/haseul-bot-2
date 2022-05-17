package notifications

import (
	"fmt"
	"log"
	"strings"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var notiGlobRemoveCommand = &router.Command{
	Name:      "remove",
	UseTyping: true,
	Run:       notiGlobRemoveRun,
}

func notiGlobRemoveRun(ctx router.CommandCtx, args []string) {
	if len(args) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a keyword to remove a notification for.",
		)
		return
	}

	go ctx.State.DeleteMessage(ctx.Msg.ChannelID, ctx.Msg.ID,
		"User removed keyword notification",
	)

	rawKeyword := util.TrimArgs(ctx.Msg.Content, ctx.Length)
	keyword := strings.ToLower(rawKeyword)

	ok, err := db.Notifications.RemoveGlobal(keyword, ctx.Msg.Author.ID)
	if err != nil {
		log.Println(err)
		dctools.SendError(ctx.State, ctx.Msg.ChannelID,
			"Error occurred while removing keyword from the database.",
		)
		return
	}
	if !ok {
		dctools.SendWarning(ctx.State, ctx.Msg.ChannelID,
			"You are already not notified of this keyword.",
		)
		return
	}

	dmChannel, err := ctx.State.CreatePrivateChannel(ctx.Msg.Author.ID)
	if err != nil {
		log.Println(err)
		dctools.SendError(ctx.State, ctx.Msg.ChannelID,
			"Error occurred while trying to DM you.",
		)
		return
	}

	dmMsg := fmt.Sprintf(
		"You will no longer be notified when '%s' is mentioned globally",
		keyword,
	)

	ctx.State.SendMessage(dmChannel.ID, dmMsg)

	dctools.SendSuccess(ctx.State, ctx.Msg.ChannelID,
		"Notification was removed successfully.",
	)
}
