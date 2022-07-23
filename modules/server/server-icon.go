package server

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var serverIconCommand = &router.SubCommand{
	Name:        "icon",
	Description: "Displays the Discord server's icon",
	Handler: &router.CommandHandler{
		Executor: serverIconExec,
	},
}

func serverIconExec(ctx router.CommandCtx) {
	guild, err := ctx.State.Guild(ctx.Interaction.GuildID)
	if dctools.ErrMissingAccess(err) {
		ctx.RespondWarning(
			"I cannot access this server.")
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while fetching server data.",
		)
		return
	}

	if guild.Icon == "" {
		ctx.RespondWarning("This server has no icon.")
		return
	}

	name := util.Possessive(guild.Name)
	title := name + " Icon"
	url := dctools.ResizeImage(guild.IconURL(), 4096)

	embed := cmdutil.ImageInfoEmbed(title, url, dctools.EmbedBackColour)

	ctx.RespondEmbed(*embed)
}
