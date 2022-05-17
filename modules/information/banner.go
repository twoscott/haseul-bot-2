package information

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var bannerCommand = &router.Command{
	Name:      "banner",
	UseTyping: true,
	Run:       bannerRun,
}

func bannerRun(ctx router.CommandCtx, args []string) {
	var userID discord.UserID
	if len(args) < 1 {
		userID = ctx.Msg.Author.ID
	} else {
		userID = dctools.ParseUserID(args[0])
	}
	if !userID.IsValid() {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Invalid user or user ID provided.",
		)
		return
	}

	user, err := ctx.State.User(userID)
	if dctools.ErrUnknownUser(err) {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg, "User does not exist.")
		return
	}
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while fetching user data.",
		)
		return
	}

	if user.Banner == "" {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			user.Username+" has no banner.",
		)
		return
	}

	name := util.Possessive(user.Username)
	title := name + " Banner"
	url := dctools.ResizeImage(user.BannerURL(), 4096)

	colour := user.Accent
	if colour == 0x000000 {
		colour, _ = ctx.State.MemberColor(ctx.Msg.GuildID, user.ID)
	}
	embed := cmdutil.ImageInfoEmbed(title, url, colour)

	dctools.EmbedReplyNoPing(ctx.State, ctx.Msg, *embed)
}
