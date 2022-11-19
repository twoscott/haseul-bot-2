package lastfm

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/shkh/lastfm-go/lastfm"
	"github.com/twoscott/haseul-bot-2/router"
)

var lastFmSetCommand = &router.SubCommand{
	Name:        "set",
	Description: "Links a Last.fm username to your Discord account.",
	Handler: &router.CommandHandler{
		Executor: lastFmSetExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:  "username",
			Description: "The Last.fm username to link to your Discord account",
			Required:    true,
		},
	},
}

func lastFmSetExec(ctx router.CommandCtx) {
	lfUser := ctx.Options.Find("username").String()
	if len(lfUser) < 1 {
		ctx.RespondError(
			"You must provide a Last.fm username to " +
				"link to your Discord account.",
		)
	}

	res, err := lf.User.GetInfo(lastfm.P{"user": lfUser})
	if err != nil {
		log.Println(err)
		lfErr := getLfError(err)

		switch lfErr.Code {
		case 6:
			ctx.RespondWarning(
				fmt.Sprintf("Last.fm user %s does not exist.", lfUser),
			)
		case 8:
			ctx.RespondError(
				fmt.Sprintf(
					"Error occurred while checking if %s is a valid username.",
					lfUser,
				),
			)
		default:
			ctx.RespondError("Unknown Last.fm Error occurred.")
		}

		return
	}

	if res.Name != "" {
		lfUser = res.Name
	}

	err = db.LastFM.SetUser(ctx.Interaction.SenderID(), lfUser)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while trying to set your Last.fm username",
		)
		return
	}

	ctx.RespondSuccess(fmt.Sprintf(
		"Your Last.fm username was set to %s. "+
			"You can now use Last.fm commands freely!",
		lfUser,
	),
	)
}
