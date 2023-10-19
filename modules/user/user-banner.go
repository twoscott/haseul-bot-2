package user

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var userBannerCommand = &router.SubCommand{
	Name:        "banner",
	Description: "Displays a Discord user's banner",
	Handler: &router.CommandHandler{
		Executor: userBannerExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.UserOption{
			OptionName:  "user",
			Description: "The user to display information for",
		},
	},
}

func userBannerExec(ctx router.CommandCtx) {
	userSnowflake, _ := ctx.Options.Find("user").SnowflakeValue()

	userID := discord.UserID(userSnowflake)
	if !userID.IsValid() {
		userID = ctx.Interaction.SenderID()
	}

	user, err := ctx.State.User(userID)
	if dctools.ErrUnknownUser(err) {
		ctx.RespondWarning("User does not exist.")
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching user data.")
		return
	}

	if user.Banner == "" {
		ctx.RespondWarning(user.DisplayOrUsername() + " has no banner.")
		return
	}

	name := util.Possessive(user.DisplayOrUsername())
	title := name + " Banner"
	url := dctools.ResizeImage(user.BannerURL(), 4096)
	embed := cmdutil.ImageInfoEmbed(title, url, user.Accent)

	ctx.RespondEmbed(*embed)
}
