package user

import (
	"fmt"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var repStreaksLeaderboardCommand = &router.SubCommand{
	Name:        "leaderboard",
	Description: "Displays a list of users with the longest rep streaks",
	Handler: &router.CommandHandler{
		Executor: repStreaksLeaderboardExec,
		Defer:    true,
	},
	Options: []discord.CommandOptionValue{
		&discord.IntegerOption{
			OptionName:  "scope",
			Description: "Where to fetch user levels from",
			Choices: []discord.IntegerChoice{
				{Name: "Server", Value: serverScope},
				{Name: "Global", Value: globalScope},
			},
		},
		&discord.IntegerOption{
			OptionName:  "users",
			Description: "The amount of top users to list",
			Min:         option.NewInt(1),
			Max:         option.NewInt(100),
		},
	},
}

func repStreaksLeaderboardExec(ctx router.CommandCtx) {
	limit, _ := ctx.Options.Find("users").IntValue()
	if limit == 0 {
		limit = 10
	}

	_, err := db.Reps.UpdateRepStreaks()
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while updating rep streaks")
		return
	}

	streaks, err := db.Reps.GetTopStreaks(limit)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching top streaks.")
		return
	}

	if len(streaks) < 1 {
		ctx.RespondWarning(
			"There are no no ongoing rep streaks to display.",
		)
		return
	}

	streakList := make([]string, 0, len(streaks))
	for i, s := range streaks {
		days := time.Since(s.FirstRep) / humanize.Day

		var uname1, uname2 string
		user1, err := ctx.State.User(s.UserID1)
		if err != nil {
			log.Println(err)
			uname1 = s.UserID1.Mention()
		} else {
			uname1 = user1.Username
		}

		user2, err := ctx.State.User(s.UserID2)
		if err != nil {
			log.Println(err)
			uname2 = s.UserID2.Mention()
		} else {
			uname2 = user2.Username
		}

		row := fmt.Sprintf(
			"%d. %s & %s (%s days)",
			i+1,
			uname1,
			uname2,
			humanize.Comma(int64(days)),
		)
		streakList = append(streakList, row)
	}

	totalStreaks, _ := db.Reps.GetTotalStreaks()
	descriptionPages := util.PagedLines(streakList, 2048, 25)
	footer := util.PluraliseWithCount("Ongoing Streak", int64(totalStreaks))

	pages := make([]router.MessagePage, len(descriptionPages))
	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		pages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Title:       "Global Streaks Leaderboard",
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
