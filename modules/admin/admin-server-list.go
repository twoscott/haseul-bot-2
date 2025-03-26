package admin

import (
	"fmt"
	"log"
	"slices"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var adminServerList = &router.SubCommand{
	Name:        "list",
	Description: "Lists servers the bot is in",
	Handler: &router.CommandHandler{
		Executor: adminServerListExec,
		Defer:    true,
	},
}

func adminServerListExec(ctx router.CommandCtx) {
	guilds, err := ctx.State.AllGuildsWithCounts()
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching servers.")
		return
	}

	slices.SortFunc(guilds, func(a, b discord.Guild) int {
		return int(b.ApproximateMembers) - int(a.ApproximateMembers)
	})

	lines := make([]string, len(guilds))
	for i, g := range guilds {
		lines[i] = fmt.Sprintf(
			"- %s (%d) - %s members",
			dctools.Bold(g.Name),
			g.ID,
			humanize.Comma(int64(g.ApproximateMembers)),
		)
	}

	descriptionPages := util.PagedLines(lines, 2048, 10)
	pages := make([]router.MessagePage, len(descriptionPages))
	footer := util.PluraliseWithCount("Server", int64(len(guilds)))

	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		pages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Title:       "Servers",
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
