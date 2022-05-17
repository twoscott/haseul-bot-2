package notifications

import (
	"fmt"
	"log"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var notiGlobAddCommand = &router.Command{
	Name:      "add",
	UseTyping: true,
	Run:       notiGlobAddRun,
}

func notiGlobAddRun(ctx router.CommandCtx, args []string) {
	if len(args) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a keyword to get notified for.",
		)
		return
	}

	go ctx.State.DeleteMessage(ctx.Msg.ChannelID, ctx.Msg.ID,
		"User added keyword notification",
	)

	notifications, err := db.Notifications.GetByUser(ctx.Msg.Author.ID)
	if err != nil {
		log.Println(err)
		dctools.SendError(ctx.State, ctx.Msg.ChannelID,
			"Error occurred while fetching Notifications from the database.",
		)
		return
	}
	if len(notifications) >= 25 {
		dctools.SendWarning(ctx.State, ctx.Msg.ChannelID,
			"You cannot have more than 25 global notifications set up.",
		)
		return
	}

	keyword, keyType := getKeyword(ctx, args)
	if len([]rune(keyword)) > 128 {
		dctools.SendWarning(ctx.State, ctx.Msg.ChannelID,
			"Keywords must be less than 128 characters in length.",
		)
		return
	}

	if len([]rune(keyword)) > 128 {
		dctools.SendWarning(ctx.State, ctx.Msg.ChannelID,
			"Keywords must be less than 128 characters in length.",
		)
		return
	}

	ok, err := db.Notifications.AddGlobal(
		keyword, ctx.Msg.Author.ID, keyType,
	)
	if err != nil {
		log.Println(err)
		dctools.SendError(ctx.State, ctx.Msg.ChannelID,
			"Error occurred while adding keyword to the database.",
		)
		return
	}
	if !ok {
		dctools.SendWarning(ctx.State, ctx.Msg.ChannelID,
			"You are already notified of this keyword.",
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
		"You will now be notified when '%s' is mentioned globally.",
		keyword,
	)

	_, err = ctx.State.SendMessage(dmChannel.ID, dmMsg)
	if dctools.ErrCannotDM(err) {
		dctools.SendWarning(ctx.State, ctx.Msg.ChannelID,
			"I am unable to DM you. "+
				"Please open your DMs to server members in your settings.",
		)
		db.Notifications.RemoveGlobal(keyword, ctx.Msg.Author.ID)
		return
	}

	dctools.SendSuccess(ctx.State, ctx.Msg.ChannelID,
		"Notification was added successfully.",
	)
}
