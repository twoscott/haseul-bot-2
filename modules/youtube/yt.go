package youtube

import (
	"fmt"
	"log"
	"strings"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	ytutil "github.com/twoscott/haseul-bot-2/utils/youtubeutil"
)

var ytCommand = &router.Command{
	Name:        "yt",
	Aliases:     []string{"youtube"},
	UseTyping:   true,
	Run:         yt,
	SubCommands: make(router.CommandMap),
}

func yt(ctx router.CommandCtx, args []string) {
	searchQuery := strings.Join(args, " ")
	videoLinks, err := ytutil.MultiSearch(searchQuery)
	if err == ytutil.ErrNoResultsFound {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			fmt.Sprintf("No results found for '%s'.", searchQuery),
		)
		return
	}
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			fmt.Sprintf(
				"Error occurred trying to fetch YouTube results for '%s'.",
				searchQuery,
			),
		)
		return
	}

	messagePages := make([]router.MessagePage, len(videoLinks))
	for i, link := range videoLinks {
		messagePages[i] = router.MessagePage{
			Content: fmt.Sprintf("%d. %s", i+1, link),
		}
	}

	cmdutil.ReplyWithConfirmationPaging(ctx, ctx.Msg, messagePages)
}
