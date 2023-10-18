package logs

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
)

var logsWelcomeTitleCommand = &router.SubCommand{
	Name:        "title",
	Description: "Edit the new member welcome title",
	Handler: &router.CommandHandler{
		Executor: logsWelcomeTitleExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:  "title",
			Description: "The title of the welcome message for new members",
			Required:    true,
			MaxLength:   option.NewInt(32),
		},
	},
}

func logsWelcomeTitleExec(ctx router.CommandCtx) {
	title := ctx.Options.Find("title").String()

	set, err := db.Guilds.SetWelcomeTitle(ctx.Interaction.GuildID, title)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while setting welcome title.")
		return
	}

	if !set {
		err := fmt.Errorf(
			"welcome title wasn't updated for %d",
			ctx.Interaction.GuildID,
		)
		log.Println(err)
		ctx.RespondError("Error occurred while setting welcome title.")
		return
	}

	ctx.RespondSuccess("Welcome title edited.")
}
