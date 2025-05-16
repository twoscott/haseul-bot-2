package cmdutil

import (
	"fmt"
	"net/http"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/httputil"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

// ImageInfoEmbedWithColour returns an embed displaying the provided image url
// along with auxiliary information.
func ImageInfoEmbedWithColour(title, url string, colour discord.Color) *discord.Embed {
	var (
		format   string = "Unknown"
		sizeMB   float64
		width    uint32
		height   uint32
		modified time.Time
		res      *http.Response
		img      *util.RawImage
	)

	res, err := util.GetImageURL(dctools.ResizeImage(url, 256))
	if err == nil {
		img, err = util.RawImageFromResponse(*res)
	}
	if err == nil {
		size := img.Size()
		sizeMB = float64(size) / humanize.MByte
		format = img.Type().String()

		dims := img.Dimensions()
		width = dims[0]
		height = dims[1]

		modified, _ = httputil.HeaderModifiedTime(res.Header)

		if colour == discord.NullColor {
			c, err := img.Colour()
			if err == nil {
				colour = dctools.ConvertColour(c)
			}
		}
	}

	embed := discord.Embed{
		Title: title,
		Image: &discord.EmbedImage{
			URL: url,
		},
		Footer: &discord.EmbedFooter{
			Text: dctools.SeparateEmbedFooter(
				fmt.Sprintf("Type: %s", format),
				fmt.Sprintf(
					"Size: %dx%d - %1.2fMB",
					width, height, sizeMB),
			),
		},
		Color: dctools.EmbedColour(colour),
	}

	if !modified.IsZero() {
		embed.Timestamp = discord.Timestamp(modified)
	}

	return &embed
}

// ImageInfoEmbedWithColour returns an embed displaying the provided image url
// along with auxiliary information. Auto-detects a colour from the image.
func ImageInfoEmbed(title, url string) *discord.Embed {
	return ImageInfoEmbedWithColour(title, url, discord.NullColor)
}

// ServerInfoEmbed returns an embed for displaying information about a guild.
func ServerInfoEmbed(st *state.State, guild discord.Guild) discord.Embed {
	url := guild.IconURLWithType(discord.PNGImage)
	colour, _ := dctools.EmbedImageColour(url)

	embed := discord.Embed{
		Title: guild.Name,
		Thumbnail: &discord.EmbedThumbnail{
			URL: guild.IconURL(),
		},
		Fields: []discord.EmbedField{},
		Color:  dctools.EmbedColour(colour),
	}

	if dctools.GuildHasFeature(guild, discord.Discoverable) {
		embed.Description = guild.Description
	}

	if guild.Banner != "" {
		url := dctools.ResizeImage(guild.BannerURL(), 4096)
		embed.Image = &discord.EmbedImage{URL: url}
	}

	var ownerValue string
	owner, err := st.User(guild.OwnerID)
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

	channels, err := st.Channels(guild.ID)
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
		Value:  dctools.Timestamp(guild.CreatedAt()),
		Inline: true,
	})

	var iconUploaded time.Time
	if guild.Icon != "" {
		iconUploaded, _ = httputil.ImgUploadTime(guild.IconURL())
	}
	if !iconUploaded.IsZero() {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:   "Icon Uploaded",
			Value:  dctools.Timestamp(iconUploaded),
			Inline: true,
		})
	}

	if guild.VanityURLCode != "" {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  "Vanity Invite URL",
			Value: fmt.Sprintf("https://discord.gg/%s", guild.VanityURLCode),
		})
	}

	return embed
}
