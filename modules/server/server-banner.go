package server

import (
	"log"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var serverBannerCommand = &router.SubCommand{
	Name:        "banner",
	Description: "Displays the Discord server's banner",
	Handler: &router.CommandHandler{
		Executor: serverBannerExec,
	},
}

func serverBannerExec(ctx router.CommandCtx) {
	guild, err := ctx.State.GuildWithCount(ctx.Interaction.GuildID)
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

	if guild.Banner == "" {
		ctx.RespondWarning(
			"This server has no banner.",
		)
		return
	}

	name := util.Possessive(guild.Name)
	title := name + " Banner"
	url := dctools.ResizeImage(guild.BannerURL(), 4096)

	embed := cmdutil.ImageInfoEmbed(title, url, dctools.EmbedBackColour)

	ctx.RespondEmbed(*embed)
}
