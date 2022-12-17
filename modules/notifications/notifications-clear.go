package notifications

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var notificationsClearCommand = &router.SubCommand{
	Name:        "clear",
	Description: "Clears all keyword notifications",
	Handler: &router.CommandHandler{
		Executor:  notificationsClearExec,
		Ephemeral: true,
	},
	Options: []discord.CommandOptionValue{
		&discord.IntegerOption{
			OptionName:  "scope",
			Description: "Where to clear keyword notifications",
			Choices: []discord.IntegerChoice{
				{Name: "Server", Value: int(serverScope)},
				{Name: "Global", Value: int(globalScope)},
			},
		},
	},
}

func notificationsClearExec(ctx router.CommandCtx) {
	scope, _ := ctx.Options.Find("scope").IntValue()

	switch scope {
	case serverScope:
		clearServerNotis(ctx)
	case globalScope:
		clearGlobalNotis(ctx)
	default:
		ctx.RespondError("Invalid notification scope selected.")
	}
}

func clearServerNotis(ctx router.CommandCtx) {
	cleared, err := db.Notifications.Clear(
		ctx.Interaction.SenderID(), ctx.Interaction.GuildID,
	)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while clearing all notifications from " +
				"the database.",
		)
		return
	}
	if cleared == 0 {
		ctx.RespondWarning(
			"You have no notifications to be cleared in this server.",
		)
		return
	}

	ctx.RespondSuccess(
		"Your notifications have been cleared from this server.",
	)
}

func clearGlobalNotis(ctx router.CommandCtx) {
	cleared, err := db.Notifications.ClearGlobal(ctx.Interaction.SenderID())
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while clearing all notifications from " +
				"the database.",
		)
		return
	}
	if cleared == 0 {
		ctx.RespondWarning(
			"You have no global notifications to be cleared.",
		)
		return
	}

	ctx.RespondSuccess(
		"Your global notifications have been cleared.",
	)
}
