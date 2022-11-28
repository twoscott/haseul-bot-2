package notification

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
	"golang.org/x/exp/slices"
)

var notificationListCommand = &router.SubCommand{
	Name:        "list",
	Description: "Lists all notifications",
	Handler: &router.CommandHandler{
		Executor:  notificationListExec,
		Ephemeral: true,
	},
}

func notificationListExec(ctx router.CommandCtx) {
	notifications, err := db.Notifications.GetByUser(ctx.Interaction.SenderID())
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while fetching Notifications from the database.",
		)
		return
	}
	if len(notifications) < 1 {
		ctx.RespondWarning(
			"You have no notifications set up with Haseul Bot.",
		)
		return
	}

	notiList := make([]string, 0, len(notifications))
	for _, noti := range notifications {
		scope := "Global"
		if noti.GuildID.IsValid() {
			g, err := ctx.State.Guild(noti.GuildID)
			if err != nil {
				scope = "Server"
			} else {
				scope = g.Name + " Server"
			}
		}

		entry := fmt.Sprintf("**%s** - %s (%s)", noti.Keyword, noti.Type, scope)
		notiList = append(notiList, entry)
	}

	slices.Sort(notiList)

	descriptionPages := util.PagedLines(notiList, 2048, 10)
	pages := make([]router.MessagePage, len(descriptionPages))
	footer := util.PluraliseWithCount("Notification", int64(len(notifications)))

	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		pages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Title:       "Notification List",
					Description: description,
					Color:       dctools.EmbedBackColour,
					Footer: &discord.EmbedFooter{
						Text: dctools.SeparateEmbedFooter(
							pageID,
							footer,
						),
					},
				},
			},
		}
	}

	ctx.RespondPaging(pages)
}
