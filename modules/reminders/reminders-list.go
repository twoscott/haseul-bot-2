package reminders

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/database/reminderdb"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
	"golang.org/x/exp/slices"
)

var remindersListCommand = &router.SubCommand{
	Name:        "list",
	Description: "List all the pending reminders you have set",
	Handler: &router.CommandHandler{
		Executor:  remindersListExec,
		Ephemeral: true,
	},
}

func remindersListExec(ctx router.CommandCtx) {
	reminders, err := db.Reminders.GetAllByUser(ctx.Interaction.SenderID())
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching reminders.")
		return
	}

	if len(reminders) < 1 {
		ctx.RespondWarning("You don't have any pending reminders.")
		return
	}

	slices.SortFunc(reminders, func(a, b reminderdb.Reminder) bool {
		return a.Time.Unix() > b.Time.Unix()
	})

	lines := make([]string, len(reminders))
	for i, r := range reminders {
		lines[i] = fmt.Sprintf(
			"%s - %s",
			dctools.UnixTimestamp(r.Time),
			r.Content,
		)
	}

	descriptionPages := util.PagedLines(lines, 2048, 10)
	pages := make([]router.MessagePage, len(descriptionPages))
	footer := util.PluraliseWithCount("Reminder", int64(len(reminders)))

	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		pages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Title:       "Pending Reminders",
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
