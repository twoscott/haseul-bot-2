package emoji

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/httputil"
)

var emojiExpandCommand = &router.SubCommand{
	Name:        "expand",
	Description: "Expands a custom Discord emoji and displays it",
	Handler: &router.CommandHandler{
		Executor: emojiExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:  "emoji",
			Description: "The Discord emoji to expand",
			Required:    true,
		},
	},
}

func emojiExec(ctx router.CommandCtx) {
	emojiString := ctx.Options.Find("emoji").String()

	emoji, err := dctools.ParseEmoji(emojiString)
	if err != nil {
		ctx.RespondWarning("Please provide a custom Discord emoji.")
		return
	}

	res, err := httputil.Head(emoji.EmojiURL())
	if res.StatusCode == http.StatusNotFound {
		ctx.RespondWarning("This emoji does not exist.")
		return
	}
	if err != nil || res.StatusCode != http.StatusOK {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while fetching emoji data.",
		)
		return
	}

	length := res.Header.Get("Content-Length")
	size, _ := strconv.ParseFloat(length, 64)
	sizeKB := size / humanize.KByte
	sizeKB = math.Floor(sizeKB*100) / 100

	embed := discord.Embed{
		Title: "`:" + emoji.Name + ":`",
		Image: &discord.EmbedImage{
			URL: emoji.EmojiURL(),
		},
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("Size: %sKB", humanize.Commaf(sizeKB)),
		},
		Timestamp: discord.Timestamp(emoji.CreatedAt()),
	}

	ctx.RespondEmbed(embed)
}
