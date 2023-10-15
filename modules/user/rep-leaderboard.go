package user

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/haseul-bot-2/database/repdb"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
	"golang.org/x/exp/slices"
)

var repLeaderboardCommand = &router.SubCommand{
	Name:        "leaderboard",
	Description: "Displays a list of users with the highest rep score",
	Handler: &router.CommandHandler{
		Executor: repLeaderboardExec,
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

func repLeaderboardExec(ctx router.CommandCtx) {
	scope, _ := ctx.Options.Find("scope").IntValue()
	limit, _ := ctx.Options.Find("users").IntValue()
	if limit == 0 {
		limit = 10
	}

	var (
		userReps []repdb.RepUser
		err      error
		listName string
	)
	switch scope {
	case serverScope:
		guild, err := ctx.State.Guild(ctx.Interaction.GuildID)
		if err != nil {
			break
		}

		members, err := ctx.State.Members(guild.ID)
		if err != nil {
			break
		}

		allReps, err := db.Reps.GetAllUsers()
		if err != nil {
			break
		}

		for _, r := range allReps {
			match := slices.ContainsFunc(members, func(m discord.Member) bool {
				return r.UserID == m.User.ID
			})

			if match {
				userReps = append(userReps, r)
			}
		}

		listName = guild.Name + " "
	case globalScope:
		listName = "Global" + " "
		userReps, err = db.Reps.GetTopUsers(limit)
	}
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching top users.")
		return
	}

	userList := make([]string, 0, len(userReps))
	for i, u := range userReps {
		var username string
		user, err := ctx.State.User(u.UserID)
		if err != nil {
			log.Println(err)
			username = u.UserID.Mention()
		} else {
			username = user.Username
		}

		row := fmt.Sprintf(
			"%d. %s (%s rep)",
			i+1,
			username,
			humanize.Comma(int64(u.Rep)),
		)
		userList = append(userList, row)
	}

	totalReps, _ := db.Reps.GetTotalReps()
	descriptionPages := util.PagedLines(userList, 2048, 25)
	footer := util.PluraliseWithCount("Total Rep", int64(totalReps))

	pages := make([]router.MessagePage, len(descriptionPages))
	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		pages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Title:       listName + "Rep Leaderboard",
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
