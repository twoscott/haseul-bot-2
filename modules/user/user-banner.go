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
			Description: "The user to display the banner for",
		},
		&discord.IntegerOption{
			OptionName:  "type",
			Description: "The type of banner to show",
			Choices: []discord.IntegerChoice{
				{Name: "Server", Value: int(serverType)},
				{Name: "Global", Value: int(globalType)},
			},
		},
	},
}

func userBannerExec(ctx router.CommandCtx) {
	userSnowflake, _ := ctx.Options.Find("user").SnowflakeValue()
	bannerType, _ := ctx.Options.Find("type").IntValue()

	guildID := ctx.Interaction.GuildID

	var member *discord.Member
	userID := discord.UserID(userSnowflake)
	if !userID.IsValid() {
		userID = ctx.Interaction.SenderID()
		member = ctx.Interaction.Member
	}
	if member == nil {
		member, _ = ctx.State.Member(guildID, userID)
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

	name := util.Possessive(user.DisplayOrUsername())
	title := name + " Banner"

	var embed *discord.Embed
	if member == nil || bannerType == globalType {
		if user.Banner == "" {
			ctx.RespondWarning(user.DisplayOrUsername() + " has no banner.")
			return
		}

		url := dctools.ResizeImage(user.BannerURL(), 4096)
		embed = cmdutil.ImageInfoEmbed(title, url)
	} else {
		url := dctools.MemberBannerURL(*member, guildID)
		if url == "" {
			ctx.RespondWarning(user.DisplayOrUsername() + " has no banner.")
			return
		}
		url = dctools.ResizeImage(url, 4096)

		colour, _ := ctx.State.MemberColor(guildID, member.User.ID)
		if dctools.ColourInvalid(colour) {
			colour = discord.NullColor
		}

		embed = cmdutil.ImageInfoEmbedWithColour(title, url, colour)
	}

	ctx.RespondEmbed(*embed)
}
