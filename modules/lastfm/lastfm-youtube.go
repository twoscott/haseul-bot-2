package lastfm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/ytutil"
)

var lastFMYouTubeCommand = &router.SubCommand{
	Name:        "youtube",
	Description: "Searches YouTube for your currently scrobbling track",
	Handler: &router.CommandHandler{
		Executor: lastFMYouTubeExec,
		Defer:    true,
	},
}

func lastFMYouTubeExec(ctx router.CommandCtx) {
	fmUser, err := db.LastFM.GetUser(ctx.Interaction.SenderID())
	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondWarning(
			fmt.Sprintf(
				"Please link a Last.fm username to your account using %s",
				lastFMSetCommand.Mention(),
			),
		)
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondGenericError()
		return
	}

	res, err := fm.User.RecentTrack(fmUser)
	if err != nil {
		log.Println(err)
		if msg, ok := errMessage(err); ok {
			ctx.RespondError(msg)
		} else {
			ctx.RespondError("Unable to fetch your recent scrobble from Last.fm.")
		}
		return
	}

	if res.Track == nil {
		ctx.RespondWarning("You have not scrobbled any tracks on Last.fm.")
		return
	}

	track := res.Track

	searchQuery := fmt.Sprintf("%s - %s", track.Artist.Name, track.Title)
	videoLinks, err := ytutil.MultiSearch(searchQuery, 20)
	if err != nil {
		if errors.Is(err, ytutil.ErrNoResultsFound) {
			ctx.RespondWarning(
				fmt.Sprintf("No results found for '%s'.", searchQuery),
			)
		} else {
			log.Println(err)
			ctx.RespondError(
				fmt.Sprintf(
					"Error occurred while fetching YouTube results for '%s'.",
					searchQuery,
				),
			)
		}
		return
	}

	var prefix string
	if track.NowPlaying {
		prefix = "Now Scrobbling"
	} else {
		prefix = "Last Scrobbled"
	}

	if len(videoLinks) == 1 {
		ctx.RespondText(fmt.Sprintf("%s: %s", prefix, videoLinks[0]))
		return
	}

	messagePages := make([]router.MessagePage, len(videoLinks))
	for i, link := range videoLinks {
		messagePages[i] = router.MessagePage{
			Content: fmt.Sprintf("%s: %s", prefix, link),
		}
	}

	ctx.RespondConfirmationPaging(messagePages)
}
