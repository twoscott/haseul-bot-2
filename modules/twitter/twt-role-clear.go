package twitter

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var twtRoleClearCommand = &router.Command{
	Name:                "clear",
	RequiredPermissions: discord.PermissionManageChannels,
	IncludeAdmin:        true,
	UseTyping:           true,
	Run:                 twtRoleClearRun,
}

func twtRoleClearRun(ctx router.CommandCtx, _ []string) {
	cleared, err := db.Twitter.ClearGuildMentions(ctx.Msg.GuildID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while clearing all Twitter roles "+
				"from the database.",
		)
		return
	}
	if cleared == 0 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"There are no Twitter roles to be cleared in this server.",
		)
		return
	}

	dctools.ReplyWithSuccess(ctx.State, ctx.Msg,
		"Twitter roles have been cleared from this server.",
	)
}
