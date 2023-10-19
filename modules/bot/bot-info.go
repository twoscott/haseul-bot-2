package bot

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/botutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var botInfoCommand = &router.SubCommand{
	Name:        "info",
	Description: "Displays information about the bot",
	Handler: &router.CommandHandler{
		Executor: botInfoExec,
	},
}

func botInfoExec(ctx router.CommandCtx) {
	bot, err := ctx.State.Me()
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching my data.")
		return
	}

	colour := bot.Accent
	if colour == 0x000000 {
		colour, _ = ctx.State.MemberColor(ctx.Interaction.GuildID, bot.ID)
	}

	embed := &discord.Embed{
		Title: bot.DisplayOrUsername() + " Info",
		Thumbnail: &discord.EmbedThumbnail{
			URL: bot.AvatarURL(),
		},
		Description: bot.Mention(),
		Fields:      []discord.EmbedField{},
		Color:       dctools.EmbedColour(colour),
	}

	var authorValue string
	author, err := ctx.State.User(botutil.AuthorID())
	if err != nil {
		authorValue = botutil.AuthorID().Mention()
	} else {
		authorValue = author.Tag()
	}

	guilds, _ := ctx.State.Guilds()
	embed.Fields = append(embed.Fields, discord.EmbedField{
		Name:  "Author",
		Value: authorValue,
	})

	members := 0
	for _, g := range guilds {
		ms, _ := ctx.State.MemberStore.Members(g.ID)
		members += len(ms)
	}

	embed.Fields = append(embed.Fields,
		discord.EmbedField{
			Name:   "Uptime",
			Value:  botutil.Uptime().String(),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Servers",
			Value:  humanize.Comma(int64(len(guilds))),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Cached Members",
			Value:  humanize.Comma(int64(members)),
			Inline: true,
		},
	)

	embed.Fields = append(embed.Fields,
		discord.EmbedField{
			Name:   "OS",
			Value:  strings.Title(runtime.GOOS),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Architecture",
			Value:  strings.Title(runtime.GOARCH),
			Inline: true,
		},
	)

	embed.Fields = append(embed.Fields, discord.EmbedField{
		Name:  "Bot Created",
		Value: dctools.UnixTimestamp(bot.CreatedAt()),
	})

	botMember, err := ctx.State.Member(ctx.Interaction.GuildID, bot.ID)
	if botMember != nil && botMember.Joined.IsValid() {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  "Bot Joined",
			Value: dctools.UnixTimestamp(botMember.Joined.Time()),
		})
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	heapMB := float64(memStats.HeapAlloc) / humanize.MByte

	site := dctools.Hyperlink("Website", botutil.Website)
	invite := dctools.Hyperlink("Discord", botutil.Discord)
	patreon := dctools.Hyperlink("Patreon", botutil.Patreon)

	embed.Fields = append(embed.Fields,
		discord.EmbedField{
			Name:  "Memory Usage",
			Value: fmt.Sprintf("%1.2fMB", heapMB),
		},
		discord.EmbedField{
			Name:  "Links",
			Value: fmt.Sprintf("%s - %s - %s", site, invite, patreon),
		},
	)

	ctx.RespondEmbed(*embed)
}
