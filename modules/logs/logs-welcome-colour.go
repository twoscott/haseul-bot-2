package logs

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var logsWelcomeColourCommand = &router.SubCommand{
	Name:        "colour",
	Description: "Edit the new member welcome embed colour",
	Handler: &router.CommandHandler{
		Executor: logsWelcomeColourExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName: "colour",
			Description: "The colour of the welcome message embed for new " +
				"members",
			Required:  true,
			MaxLength: option.NewInt(7),
		},
	},
}

func logsWelcomeColourExec(ctx router.CommandCtx) {
	colourString := ctx.Options.Find("colour").String()

	colour, err := dctools.HexToColour(colourString)
	if err != nil {
		log.Println(err)
		ctx.RespondWarning("Provided hex colour value is invalid.")
		return
	}

	set, err := db.Guilds.SetWelcomeColour(ctx.Interaction.GuildID, colour)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while setting welcome colour.")
		return
	}

	if !set {
		err := fmt.Errorf(
			"welcome colour wasn't updated for %d",
			ctx.Interaction.GuildID,
		)
		log.Println(err)
		ctx.RespondError("Error occurred while setting welcome colour.")
		return
	}

	ctx.RespondSuccess("Welcome colour edited.")
}
