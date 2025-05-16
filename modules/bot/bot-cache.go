package bot

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/botutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var botCacheCommand = &router.SubCommand{
	Name:        "cache",
	Description: "Displays information the bot's memory usage",
	Handler: &router.CommandHandler{
		Executor: botCacheExec,
	},
}

func botCacheExec(ctx router.CommandCtx) {
	bot, err := ctx.State.Me()
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching my data.")
		return
	}

	embed := &discord.Embed{
		Title:  bot.DisplayOrUsername() + " Cache Stats",
		Fields: []discord.EmbedField{},
		Color:  dctools.EmbedBackColour,
	}

	guilds, _ := ctx.State.GuildStore.Guilds()
	channels := 0
	members := 0
	roles := 0
	emojis := 0
	messages := 0

	for _, g := range guilds {
		chs, _ := ctx.State.ChannelStore.Channels(g.ID)
		channels += len(chs)

		for _, ch := range chs {
			ms, _ := ctx.State.MessageStore.Messages(ch.ID)
			messages += len(ms)
		}

		ms, _ := ctx.State.MemberStore.Members(g.ID)
		members += len(ms)

		rs, _ := ctx.State.RoleStore.Roles(g.ID)
		roles += len(rs)

		es, _ := ctx.State.EmojiStore.Emojis(g.ID)
		emojis += len(es)
	}

	embed.Fields = append(embed.Fields,
		discord.EmbedField{
			Name:   "Servers",
			Value:  humanize.Comma(int64(len(guilds))),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Channels",
			Value:  humanize.Comma(int64(channels)),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Members",
			Value:  humanize.Comma(int64(members)),
			Inline: true,
		},
	)

	embed.Fields = append(embed.Fields,
		discord.EmbedField{
			Name:   "Roles",
			Value:  humanize.Comma(int64(roles)),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Emojis",
			Value:  humanize.Comma(int64(emojis)),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Messages",
			Value:  humanize.Comma(int64(messages)),
			Inline: true,
		},
	)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	heapMB := float64(memStats.HeapAlloc) / humanize.MByte
	inUseMB := float64(memStats.HeapInuse) / humanize.MByte
	idleMB := float64(memStats.HeapIdle) / humanize.MByte
	stackUseMB := float64(memStats.StackInuse) / humanize.MByte
	totalGB := float64(memStats.TotalAlloc) / humanize.GByte

	embed.Fields = append(embed.Fields,
		discord.EmbedField{
			Name:   "Heap Allocation",
			Value:  fmt.Sprintf("%1.2fMB", heapMB),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Heap In Use",
			Value:  fmt.Sprintf("%1.2fMB", inUseMB),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Heap Idle",
			Value:  fmt.Sprintf("%1.2fMB", idleMB),
			Inline: true,
		},
	)

	embed.Fields = append(embed.Fields,
		discord.EmbedField{
			Name:   "Stack In Use",
			Value:  fmt.Sprintf("%1.2fMB", stackUseMB),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Cumulative Allocation",
			Value:  fmt.Sprintf("%1.2fGB", totalGB),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Uptime",
			Value:  botutil.Uptime().String(),
			Inline: true,
		},
	)

	var gcStats debug.GCStats
	debug.ReadGCStats(&gcStats)

	embed.Fields = append(embed.Fields,
		discord.EmbedField{
			Name:   "Last GC",
			Value:  dctools.Timestamp(gcStats.LastGC),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Num Of GCs",
			Value:  humanize.Comma(gcStats.NumGC),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "GC Total Pauses",
			Value:  gcStats.PauseTotal.String(),
			Inline: true,
		},
	)

	ctx.RespondEmbed(*embed)
}
