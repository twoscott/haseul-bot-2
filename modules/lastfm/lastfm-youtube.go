package lastfm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/ytutil"
)

var lastFmYouTubeCommand = &router.SubCommand{
	Name:        "youtube",
	Description: "Searches YouTube for your currently scrobbling track",
	Handler: &router.CommandHandler{
		Executor: lastFmYouTubeExec,
		Defer:    true,
	},
}

func lastFmYouTubeExec(ctx router.CommandCtx) {
	lfUser, err := db.LastFM.GetUser(ctx.Interaction.SenderID())
	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondWarning(
			fmt.Sprintf(
				"Please link a Last.fm username to your account using %s",
				lastFmSetCommand.Mention(),
			),
		)
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondGenericError()
		return
	}

	res, err := getRecentTracks(lfUser, 1)
	if err != nil {
		errMsg := errorResponseMessage(err)
		ctx.RespondError(errMsg)
		return
	}

	if len(res.Tracks) < 1 {
		ctx.RespondWarning(
			"You have not scrobbled any tracks on Last.fm.",
		)
		return
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
	videoLinks, err := ytutil.MultiSearch(searchQuery, 20)
	if errors.Is(err, ytutil.ErrNoResultsFound) {
		ctx.RespondWarning(
			fmt.Sprintf("No results found for '%s'.", searchQuery),
		)
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			fmt.Sprintf(
				"Error occurred while fetching YouTube results for '%s'.",
				searchQuery,
			),
		)
		return
	}

	if len(videoLinks) == 1 {
		youTubeResponse := fmt.Sprintf("%s: %s", prefix, videoLinks[0])
		ctx.RespondText(youTubeResponse)
	}

	messagePages := make([]router.MessagePage, len(videoLinks))
	for i, link := range videoLinks {
		messagePages[i] = router.MessagePage{
			Content: fmt.Sprintf("%s: %s", prefix, link),
		}
	}

	ctx.RespondConfirmationPaging(messagePages)
}
