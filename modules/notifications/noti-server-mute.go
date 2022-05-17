package notifications

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var notiGuildMuteCommand = &router.Command{
	Name:      "mute",
	Aliases:   []string{"blacklist", "ignore"},
	UseTyping: true,
	Run:       notiGuildMuteRun,
}

func notiGuildMuteRun(ctx router.CommandCtx, _ []string) {
	muted, err := db.Notifications.MuteChannel(
		ctx.Msg.Author.ID, ctx.Msg.ChannelID,
	)
	if err != nil {
		log.Println(err)
		dctools.SendError(ctx.State, ctx.Msg.ChannelID,
			"Error occurred while trying to the mute the channel.",
		)
		return
	}

	if muted {
		dctools.SendSuccess(ctx.State, ctx.Msg.ChannelID,
			"You will no longer be notified for keywords "+
				"mentioned in this server.",
		)
	} else {
		dctools.SendWarning(ctx.State, ctx.Msg.ChannelID,
			"This server is already muted.",
		)
	}
}
