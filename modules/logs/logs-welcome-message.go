package logs

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
)

var logsWelcomeMessageCommand = &router.SubCommand{
	Name:        "message",
	Description: "Edit the new member welcome message",
	Handler: &router.CommandHandler{
		Executor: logsWelcomeMessageExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:  "message",
			Description: "The channel to welcome new members in",
			Required:    true,
			MaxLength:   option.NewInt(1024),
		},
	},
}

func logsWelcomeMessageExec(ctx router.CommandCtx) {
	message := ctx.Options.Find("message").String()

	set, err := db.Guilds.SetWelcomeMessage(ctx.Interaction.GuildID, message)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while setting welcome message.")
		return
	}

	if !set {
		err := fmt.Errorf(
			"welcome message wasn't updated for %d",
			ctx.Interaction.GuildID,
		)
		log.Println(err)
		ctx.RespondError("Error occurred while setting welcome message.")
		return
	}

	ctx.RespondSuccess("Welcome message edited.")
}
