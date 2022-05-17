package lastfm

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	ytutil "github.com/twoscott/haseul-bot-2/utils/youtubeutil"
)

var fmYtCommand = &router.Command{
	Name:      "yt",
	Aliases:   []string{"youtube"},
	UseTyping: true,
	Run:       fmYtRun,
}

func fmYtRun(ctx router.CommandCtx, _ []string) {
	lfUser, ok := getLfUser(ctx)
	if !ok {
		return
	}

	res, ok := recentTracks(ctx, lfUser, 1)
	if !ok {
		return
	}

	// correct extra track returned when now playing
	if len(res.Tracks) > 1 {
		res.Tracks = res.Tracks[:1]
	}

	track := res.Tracks[0]
	nowPlaying, _ := strconv.ParseBool(track.NowPlaying)

	var prefix string
	if nowPlaying {
		prefix = "Now Scrobbling"
	} else {
		prefix = "Last Scrobbled"
	}

	searchQuery := fmt.Sprintf("%s - %s", track.Artist.Name, track.Name)
	videoLinks, err := ytutil.MultiSearch(searchQuery)
	if errors.Is(err, ytutil.ErrNoResultsFound) {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			fmt.Sprintf("No results found for '%s'.", searchQuery),
		)
		return
	}
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			fmt.Sprintf(
				"Error occurred while fetching YouTube results for '%s'.",
				searchQuery,
			),
		)
		return
	}

	if len(videoLinks) == 1 {
		youTubeResponse := fmt.Sprintf("%s: %s", prefix, videoLinks[0])
		dctools.TextReplyNoPing(ctx.State, ctx.Msg, youTubeResponse)
	}

	messagePages := make([]router.MessagePage, len(videoLinks))
	for i, link := range videoLinks {
		messagePages[i] = router.MessagePage{
			Content: fmt.Sprintf("%s: %s", prefix, link),
		}
	}

	cmdutil.ReplyWithConfirmationPaging(ctx, ctx.Msg, messagePages)
}
