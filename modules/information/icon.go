package information

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/botutil"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var iconCommand = &router.Command{
	Name:      "icon",
	Aliases:   []string{"guildicon", "servericon"},
	UseTyping: true,
	Run:       iconRun,
}

func iconRun(ctx router.CommandCtx, args []string) {
	var guildID discord.GuildID
	if botutil.IsBotAdmin(ctx.Msg.Author.ID) && len(args) > 0 {
		guildID = dctools.ParseGuildID(args[0])
	} else {
		guildID = ctx.Msg.GuildID
	}
	if !guildID.IsValid() {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg, "Invalid server ID provided.")
		return
	}

	guild, err := ctx.State.GuildWithCount(guildID)
	if dctools.ErrMissingAccess(err) {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"I cannot access this server.")
		return
	}
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while fetching server data.",
		)
		return
	}

	if guild.Icon == "" {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg, "This server has no icon.")
		return
	}

	name := util.Possessive(guild.Name)
	title := name + " Icon"
	url := dctools.ResizeImage(guild.IconURL(), 4096)

	embed := cmdutil.ImageInfoEmbed(title, url, dctools.EmbedBackColour)

	dctools.EmbedReplyNoPing(ctx.State, ctx.Msg, *embed)
}
