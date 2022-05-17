package lastfm

import (
	"fmt"
	"log"

	"github.com/shkh/lastfm-go/lastfm"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var fmSetCommand = &router.Command{
	Name:      "set",
	UseTyping: true,
	Run:       fmSetRun,
}

func fmSetRun(ctx router.CommandCtx, args []string) {
	if len(args) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a Last.fm username to link to your "+
				"Discord account.",
		)
		return
	}

	lfUser := args[0]

	if len(lfUser) > 15 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Last.fm usernames must be shorter than 15 characters.",
		)
		return
	}

	res, err := lf.User.GetInfo(lastfm.P{"user": lfUser})
	if err != nil {
		lfErr := getLfError(err)
		switch lfErr.Code {
		case 6:
			dctools.ReplyWithWarning(ctx.State, ctx.Msg,
				fmt.Sprintf("Last.fm user %s does not exist.", lfUser),
			)
		case 8:
			log.Println(err)
			dctools.ReplyWithError(ctx.State, ctx.Msg,
				fmt.Sprintf(
					"Error occurred while checking if %s is a valid username.",
					lfUser,
				),
			)
		default:
			log.Println(err)
			dctools.ReplyWithError(ctx.State, ctx.Msg,
				"Unknown Last.fm Error occurred.",
			)
		}
		return
	}

	if res.Name != "" {
		lfUser = res.Name
	}

	err = db.LastFM.SetUpdateUser(ctx.Msg.Author.ID, lfUser)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while trying to set your Last.fm username",
		)
		return
	}

	dctools.ReplyWithSuccess(ctx.State, ctx.Msg,
		fmt.Sprintf(
			"Your Last.fm username was set to %s. "+
				"You can now use Last.fm commands freely!",
			lfUser,
		),
	)
}
