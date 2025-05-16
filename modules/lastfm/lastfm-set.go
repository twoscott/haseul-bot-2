package lastfm

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"

	"github.com/twoscott/gobble-fm/api"
	"github.com/twoscott/haseul-bot-2/router"
)

var lastFMSetCommand = &router.SubCommand{
	Name:        "set",
	Description: "Links a Last.fm username to your Discord account",
	Handler: &router.CommandHandler{
		Executor: lastFMSetExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:  "username",
			Description: "The Last.fm username to link to your Discord account",
			Required:    true,
		},
	},
}

func lastFMSetExec(ctx router.CommandCtx) {
	fmUser := ctx.Options.Find("username").String()

	setUser, err := db.LastFM.GetUser(ctx.Interaction.SenderID())
	if err == nil && fmUser == setUser {
		m := "You already have Last.fm username '%s' linked to your Discord account."
		ctx.RespondWarning(fmt.Sprintf(m, setUser))
		return
	}

	user, err := fm.User.Info(fmUser)
	if err != nil {
		log.Println(err)
		fmerr, ok := fmError(err)
		if ok && fmerr.Code == api.ErrInvalidParameters {
			ctx.RespondWarning(
				fmt.Sprintf("Invalid Last.fm username '%s' provided.", fmUser),
			)
		} else {
			ctx.RespondError("Unable to fetch your recent scrobbles from Last.fm.")
		}
		return
	}

	if user.Name != "" {
		fmUser = user.Name
	}

	err = db.LastFM.SetUser(ctx.Interaction.SenderID(), fmUser)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while trying to set your Last.fm username",
		)
		return
	}

	m := "Your Last.fm username was set to %s. You can now use Last.fm commands!"
	ctx.RespondSuccess(fmt.Sprintf(m, fmUser))
}
