package admin

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var adminServerInfoCommand = &router.SubCommand{
	Name:        "info",
	Description: "Displays information about a Discord server",
	Handler: &router.CommandHandler{
		Executor: adminServerInfoExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:  "server",
			Description: "Server ID to fetch info for",
			Required:    true,
		},
	},
}

func adminServerInfoExec(ctx router.CommandCtx) {
	snowflake, _ := ctx.Options.Find("server").SnowflakeValue()
	guildID := discord.GuildID(snowflake)
	if !guildID.IsValid() {
		ctx.RespondWarning("Malformed server ID provided.")
		return
	}

	guild, err := ctx.State.GuildWithCount(guildID)
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
