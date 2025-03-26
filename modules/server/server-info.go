package server

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var serverInfoCommand = &router.SubCommand{
	Name:        "info",
	Description: "Displays information about the Discord server",
	Handler: &router.CommandHandler{
		Executor: serverInfoExec,
	},
}

func serverInfoExec(ctx router.CommandCtx) {
	guild, err := ctx.State.GuildWithCount(ctx.Interaction.GuildID)
	if dctools.ErrMissingAccess(err) {
		ctx.RespondWarning("I cannot access this server.")
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching server data.")
		return
	}

	embed := cmdutil.ServerInfoEmbed(ctx.State, *guild)
	ctx.RespondEmbed(embed)
}
