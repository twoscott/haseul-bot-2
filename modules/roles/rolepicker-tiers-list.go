package roles

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var rolePickerTiersList = &router.SubCommand{
	Name:        "list",
	Description: "Lists all role tiers added to the role picker",
	Handler: &router.CommandHandler{
		Executor: rolePickerTiersListExec,
	},
}

func rolePickerTiersListExec(ctx router.CommandCtx) {
	tiers, err := db.Roles.GetAllTiersByGuild(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching role tiers.")
		return
	}
	if len(tiers) < 1 {
		ctx.RespondWarning("This server has no role tiers added to it.")
		return
	}

	tierList := make([]string, 0, len(tiers))
	for _, tier := range tiers {
		tierList = append(tierList, "- ", tier.Title())
	}

	descriptionPages := util.PagedLines(tierList, 2048, 20)
	pages := make([]router.MessagePage, len(descriptionPages))
	footer := util.PluraliseWithCount("Tier", int64(len(tierList)))

	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		pages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Name: "Role Picker Tiers",
					},
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
