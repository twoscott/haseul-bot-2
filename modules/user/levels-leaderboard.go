package user

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/haseul-bot-2/database/levelsdb"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var levelsLeaderboardCommand = &router.SubCommand{
	Name: "leaderboard",
	Description: "Lists the users with the highest levels in a " +
		"server or globally",
	Handler: &router.CommandHandler{
		Executor: levelsLeaderboardExec,
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

func levelsLeaderboardExec(ctx router.CommandCtx) {
	scope, _ := ctx.Options.Find("scope").IntValue()
	limit, _ := ctx.Options.Find("users").IntValue()
	if limit == 0 {
		limit = 10
	}

	var (
		usersXP  []levelsdb.UserXP
		err      error
		listName string
		entries  int64
	)
	switch scope {
	case serverScope:
		var gUsers []levelsdb.GuildUserXP
		gUsers, err = db.Levels.GetTopUsers(ctx.Interaction.GuildID, limit)
		for _, gu := range gUsers {
			usersXP = append(usersXP, gu.UserXP)
		}

		guild, err := ctx.State.Guild(ctx.Interaction.GuildID)
		if err != nil {
			log.Println(err)
			break
		}

		listName = guild.Name + " "
		entries, _ = db.Levels.GetEntriesSize(guild.ID)

	case globalScope:
		usersXP, err = db.Levels.GetTopGlobalUsers(limit)
		listName = "Global" + " "
		entries, _ = db.Levels.GetGlobalEntriesSize()
	}
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching top users.")
		return
	}

	userList := make([]string, 0, len(usersXP))
	for i, uxp := range usersXP {
		var username string
		user, err := ctx.State.User(uxp.UserID)
		if err != nil {
			log.Println(err)
			username = uxp.UserID.Mention()
		} else {
			username = user.Username
		}

		row := fmt.Sprintf(
			"%d. %s (Lvl %s) - %s XP",
			i+1,
			username,
			humanize.Comma(int64(uxp.Level())),
			humanize.Comma(uxp.XP),
		)
		userList = append(userList, row)
	}

	descriptionPages := util.PagedLines(userList, 2048, 25)
	footer := util.PluraliseWithCount("Total Entry", entries)

	pages := make([]router.MessagePage, len(descriptionPages))
	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		pages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Title:       listName + "Leaderboard",
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
