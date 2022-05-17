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
)

var emojiCommand = &router.Command{
	Name:        "emoji",
	Aliases:     []string{"emote"},
	UseTyping:   true,
	Run:         emojiRun,
	SubCommands: make(router.CommandMap),
}

func emojiRun(ctx router.CommandCtx, args []string) {
	if len(args) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg, "Please provide an emoji.")
		return
	}

	emoji, err := dctools.ParseEmoji(args[0])
	if err != nil {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg, "Invalid emoji provided.")
		return
	}

	res, err := http.Head(emoji.EmojiURL())
	if res.StatusCode == http.StatusNotFound {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg, "This emoji does not exist.")
		return
	}
	if err != nil || res.StatusCode != http.StatusOK {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
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

	dctools.EmbedReplyNoPing(ctx.State, ctx.Msg, embed)
}
