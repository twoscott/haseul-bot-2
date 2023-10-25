package roles

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var joinRolesList = &router.SubCommand{
	Name:        "list",
	Description: "Lists all roles set to be assigned to new members",
	Handler: &router.CommandHandler{
		Executor: joinRolesListExec,
	},
}

func joinRolesListExec(ctx router.CommandCtx) {
	// roles, err := db.Roles.GetAllRolesByTier(tier.ID)
	// if err != nil {
	// 	log.Println(err)
	// 	ctx.RespondError("Error occurred while fetching roles.")
	// 	return
	// }
	// if len(roles) < 1 {
	// 	ctx.RespondWarning("This tier has no roles added to it.")
	// 	return
	// }

	roleIDs, err := db.Roles.GetAllGuildJoinRoles(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching join roles.")
		return
	}
	if len(roleIDs) < 1 {
		ctx.RespondWarning("This server has no join roles added to it.")
		return
	}

	roleList := make([]string, 0, len(roleIDs))
	for _, id := range roleIDs {
		roleList = append(roleList, "- "+id.Mention())
	}

	name := "Server"
	guild, err := ctx.State.Guild(ctx.Interaction.GuildID)
	if err == nil {
		name = guild.Name
	}

	descriptionPages := util.PagedLines(roleList, 2048, 20)
	pages := make([]router.MessagePage, len(descriptionPages))
	footer := util.PluraliseWithCount("Role", int64(len(roleList)))

	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		pages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Name: fmt.Sprintf(
							"%s Join Roles", name,
						),
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
