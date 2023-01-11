package roles

import (
	"fmt"
	"log"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var rolePickerTiersAdd = &router.SubCommand{
	Name:        "add",
	Description: "Add a role for users to be able to select",
	Handler: &router.CommandHandler{
		Executor: rolePickerTiersAddExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:  "name",
			Description: "The name of the role tier",
			MaxLength:   option.NewInt(32),
			Required:    true,
		},
		&discord.StringOption{
			OptionName:  "description",
			Description: "The description of the role tier",
			MaxLength:   option.NewInt(1024),
		},
	},
}

func rolePickerTiersAddExec(ctx router.CommandCtx) {
	tierName := ctx.Options.Find("name").String()
	tierName = strings.ToLower(tierName)
	description := ctx.Options.Find("description").String()

	ok, err := db.RolePicker.AddTier(ctx.Interaction.GuildID, tierName, description)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred adding role tier.")
		return
	}
	if !ok {
		ctx.RespondWarning(
			"A role tier with this name already exists in the server.",
		)
		return
	}

	formattedName := util.TitleCase(tierName)
	ctx.RespondSuccess(
		fmt.Sprintf("Role tier '%s' added to the role picker.", formattedName),
	)
}
