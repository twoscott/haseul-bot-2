package user

import (
	"fmt"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var repStreaksListCommand = &router.SubCommand{
	Name:        "list",
	Description: "Displays a list of rep streaks you currently have.",
	Handler: &router.CommandHandler{
		Executor: repStreaksListExec,
		Defer:    true,
	},
}

func repStreaksListExec(ctx router.CommandCtx) {
	senderID := ctx.Interaction.SenderID()

	_, err := db.Reps.UpdateRepStreaks()
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while updating rep streaks")
		return
	}

	streaks, err := db.Reps.GetUserStreaks(senderID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching your rep streaks.")
		return
	}

	if len(streaks) < 1 {
		ctx.RespondWarning(
			"You have no ongoing rep streaks to display.",
		)
		return
	}

	streakList := make([]string, 0, len(streaks))
	for _, s := range streaks {
		otherUserID := s.OtherUser(senderID)
		days := time.Since(s.FirstRep) / humanize.Day

		var username string
		user, err := ctx.State.User(otherUserID)
		if err != nil {
			log.Println(err)
			username = otherUserID.Mention()
		} else {
			username = user.Username
		}

		row := fmt.Sprintf(
			"- %s (%s days)",
			dctools.EscapeMarkdown(username),
			humanize.Comma(int64(days)),
		)
		streakList = append(streakList, row)
	}

	descriptionPages := util.PagedLines(streakList, 2048, 25)
	footer := util.PluraliseWithCount("Ongoing Streak", int64(len(streaks)))

	pages := make([]router.MessagePage, len(descriptionPages))
	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		pages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Title:       "Ongoing Rep Streaks",
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
