package server

import (
	"fmt"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/httputil"
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

	embed := guildEmbed(&ctx, guild)
	ctx.RespondEmbed(*embed)
}

func guildEmbed(ctx *router.CommandCtx, guild *discord.Guild) *discord.Embed {
	embed := discord.Embed{
		Title: guild.Name,
		Thumbnail: &discord.EmbedThumbnail{
			URL: guild.IconURL(),
		},
		Fields: []discord.EmbedField{},
		Color:  dctools.EmbedBackColour,
	}

	if dctools.GuildHasFeature(guild, discord.Discoverable) {
		embed.Description = guild.Description
	}

	if guild.Banner != "" {
		url := dctools.ResizeImage(guild.BannerURL(), 4096)
		embed.Image = &discord.EmbedImage{URL: url}
	}

	var ownerValue string
	owner, err := ctx.State.User(guild.OwnerID)
	if err != nil {
		ownerValue = guild.OwnerID.Mention()
	} else {
		ownerValue = owner.Tag()
	}

	embed.Fields = append(embed.Fields, discord.EmbedField{
		Name:  "Owner",
		Value: ownerValue,
	})

	online, offline := dctools.GuildStatuses(guild)
	embed.Fields = append(embed.Fields, discord.EmbedField{
		Name: "Presences",
		Value: fmt.Sprintf(
			"%s%s Online %s%s Offline",
			dctools.OnlineEmoji(), humanize.Comma(int64(online)),
			dctools.OfflineEmoji(), humanize.Comma(int64(offline)),
		),
	})

	region := dctools.GuildRegionText(guild.VoiceRegion)
	levelEmoji := dctools.BoostLevelEmoji(guild.NitroBoost)
	boostersEmoji := dctools.BoostersEmoji()
	embed.Fields = append(embed.Fields,
		discord.EmbedField{
			Name:   "Region",
			Value:  region,
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Level",
			Value:  fmt.Sprintf("%s %d", levelEmoji, guild.NitroBoost),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Boosters",
			Value:  fmt.Sprintf("%s %d", boostersEmoji, guild.NitroBoosters),
			Inline: true,
		},
	)

	channels, err := ctx.State.Channels(guild.ID)
	voice := 0
	if err == nil {
		for _, ch := range channels {
			if ch.Type == discord.GuildVoice {
				voice++
			}
		}
	}
	embed.Fields = append(embed.Fields,
		discord.EmbedField{
			Name:   "Members",
			Value:  humanize.Comma(int64(guild.ApproximateMembers)),
			Inline: true,
		},
		discord.EmbedField{
			Name: "Channels",
			Value: fmt.Sprintf(
				"%s (%d Voice)",
				humanize.Comma(int64(len(channels))), voice,
			),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Roles",
			Value:  humanize.Comma(int64(len(guild.Roles))),
			Inline: true,
		},
	)

	limit := dctools.GuildEmojiLimit(guild.NitroBoost)
	static := 0
	animated := 0
	for _, e := range guild.Emojis {
		if !e.Animated {
			static++
		} else {
			animated++
		}
	}
	embed.Fields = append(embed.Fields,
		discord.EmbedField{
			Name:   "Static Emojis",
			Value:  fmt.Sprintf("%d/%d", static, limit),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Animated Emojis",
			Value:  fmt.Sprintf("%d/%d", animated, limit),
			Inline: true,
		},
	)

	embed.Fields = append(embed.Fields, discord.EmbedField{
		Name:  "Server ID",
		Value: guild.ID.String(),
	})

	streamQuality := dctools.GuildStreamQuality(guild.NitroBoost)
	audioQuality := dctools.GuildAudioQualityKbps(guild.NitroBoost)
	uploadLimit := dctools.GuildUploadLimitMB(guild.NitroBoost)
	embed.Fields = append(embed.Fields,
		discord.EmbedField{
			Name:   "Upload Limit",
			Value:  fmt.Sprintf("%dMB", uploadLimit),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Audio Quality",
			Value:  fmt.Sprintf("%dKbps", audioQuality),
			Inline: true,
		},
		discord.EmbedField{
			Name:   "Stream Quality",
			Value:  streamQuality,
			Inline: true,
		},
	)

	embed.Fields = append(embed.Fields, discord.EmbedField{
		Name:   "Server Created",
		Value:  dctools.EmbedTime(guild.CreatedAt()),
		Inline: true,
	})

	var iconUploaded time.Time
	if guild.Icon != "" {
		iconUploaded, _ = httputil.ImgUploadTime(guild.IconURL())
	}
	if !iconUploaded.IsZero() {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:   "Icon Uploaded",
			Value:  dctools.EmbedTime(iconUploaded),
			Inline: true,
		})
	}

	if guild.VanityURLCode != "" {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  "Vanity Invite URL",
			Value: fmt.Sprintf("https://discord.gg/%s", guild.VanityURLCode),
		})
	}

	return &embed
}
