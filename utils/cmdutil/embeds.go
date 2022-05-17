package cmdutil

import (
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

// ImageInfoEmbed returns an embed displaying the provided image url
// along with auxiliary information.
func ImageInfoEmbed(title, url string, colour discord.Color) *discord.Embed {
	var (
		format   string = "Unknown"
		sizeMB   float64
		width    uint32
		height   uint32
		modified time.Time
	)

	image, res, err := util.ImageFromURL(url)
	if err == nil {
		size := image.Size()
		sizeMB = float64(size) / humanize.MByte
		format = image.Type().String()

		dims := image.Dimensions()
		width = dims[0]
		height = dims[1]

		modified, _ = util.HeaderModifiedTime(res.Header)
	}

	embed := discord.Embed{
		Title: title,
		Image: &discord.EmbedImage{
			URL: url,
		},
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf(
				"Type: %s%sSize: %dx%d - %1.2fMB",
				format,
				dctools.EmbedFooterSep,
				width, height, sizeMB,
			),
		},
		Color: dctools.EmbedColour(colour),
	}

	if !modified.IsZero() {
		embed.Timestamp = discord.Timestamp(modified)
	}

	return &embed
}
