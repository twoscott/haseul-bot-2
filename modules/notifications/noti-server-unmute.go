package notifications

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var notiGuildUnmuteCommand = &router.Command{
	Name:      "unmute",
	Aliases:   []string{"whitelist", "unignore"},
	UseTyping: true,
	Run:       notiGuildUnmuteRun,
}

func notiGuildUnmuteRun(ctx router.CommandCtx, _ []string) {
	muted, err := db.Notifications.UnmuteChannel(
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
			"You will now be notified for keywords "+
				"mentioned in this server.",
		)
	} else {
		dctools.SendWarning(ctx.State, ctx.Msg.ChannelID,
			"This server is already unmuted.",
		)
	}
}
