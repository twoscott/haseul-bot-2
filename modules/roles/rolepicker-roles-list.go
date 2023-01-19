package roles

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var rolePickerRolesList = &router.SubCommand{
	Name:        "list",
	Description: "Lists all roles added to a role tier",
	Handler: &router.CommandHandler{
		Executor:      rolePickerRolesListExec,
		Autocompleter: tierNameCompleter,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:   "tier",
			Description:  "The tier to add the role to",
			MaxLength:    option.NewInt(32),
			Required:     true,
			Autocomplete: true,
		},
	},
}

func rolePickerRolesListExec(ctx router.CommandCtx) {
	tierName := ctx.Options.Find("tier").String()

	tier, err := db.Roles.GetTierByName(ctx.Interaction.GuildID, tierName)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondWarning("This role tier does not exist.")
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching role tiers.")
		return
	}

	roles, err := db.Roles.GetAllRolesByTier(tier.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching roles.")
		return
	}
	if len(roles) < 1 {
		ctx.RespondWarning("This tier has no roles added to it.")
		return
	}

	roleList := make([]string, 0, len(roles))
	for _, role := range roles {
		roleList = append(roleList, role.ID.Mention())
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
							"%s Roles", tier.Title(),
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
