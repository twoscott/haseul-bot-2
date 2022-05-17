package notifications

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var notiListCommand = &router.Command{
	Name:      "list",
	UseTyping: true,
	Run:       notiListRun,
}

func notiListRun(ctx router.CommandCtx, _ []string) {
	notifications, err := db.Notifications.GetByUser(ctx.Msg.Author.ID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while fetching Notifications from the database.",
		)
		return
	}
	if len(notifications) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"You have no notifications set up with Haseul Bot.",
		)
		return
	}

	notiList := make([]string, len(notifications))
	for i, noti := range notifications {
		scope := "Global"
		if noti.GuildID.IsValid() {
			scope = "Server"
			g, err := ctx.State.Guild(noti.GuildID)
			if err == nil {
				scope = g.Name + " Server"
			}
		}

		entry := fmt.Sprintf("'%s' - %s (%s)", noti.Keyword, noti.Type, scope)

		notiList[i] = entry
	}

	descriptionPages := util.PagedLines(notiList, 2048, 10)
	pages := make([]router.MessagePage, len(descriptionPages))
	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		pages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Title:       "Notification List",
					Description: description,
					Color:       dctools.EmbedBackColour,
					Footer:      &discord.EmbedFooter{Text: pageID},
				},
			},
		}
	}

	dmChannel, err := ctx.State.CreatePrivateChannel(ctx.Msg.Author.ID)
	if err != nil {
		log.Println(err)
		dctools.SendError(ctx.State, ctx.Msg.ChannelID,
			"Error occurred while trying to DM you.",
		)
		return
	}

	_, err = cmdutil.SendWithPaging(ctx, dmChannel.ID, pages)
	if err == nil {
		dctools.ReplyWithSuccess(ctx.State, ctx.Msg,
			"I sent a list of your notifications to your DMs.",
		)
		return
	}

	if dctools.ErrCannotDM(err) {
		dctools.SendWarning(ctx.State, ctx.Msg.ChannelID,
			"I am unable to DM you. "+
				"Please open your DMs to server members in your settings.",
		)
		return
	}

	dctools.ReplyWithError(ctx.State, ctx.Msg,
		"Error occurred while sending your list of notifications "+
			"to your DMs.",
	)
}
