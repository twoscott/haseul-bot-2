package user

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var userAvatarCommand = &router.SubCommand{
	Name:        "avatar",
	Description: "Displays a Discord user's avatar",
	Handler: &router.CommandHandler{
		Executor: userAvatarExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.UserOption{
			OptionName:  "user",
			Description: "The user to display the avatar for",
		},
		&discord.IntegerOption{
			OptionName:  "type",
			Description: "The type of avatar to show",
			Choices: []discord.IntegerChoice{
				{Name: "Server", Value: int(serverType)},
				{Name: "Global", Value: int(globalType)},
			},
		},
	},
}

func userAvatarExec(ctx router.CommandCtx) {
	userSnowflake, _ := ctx.Options.Find("user").SnowflakeValue()
	avatarType, _ := ctx.Options.Find("type").IntValue()

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
	title := name + " Avatar"

	var embed *discord.Embed
	if member == nil || avatarType == globalType {
		url := dctools.ResizeImage(user.AvatarURL(), 4096)
		embed = cmdutil.ImageInfoEmbedWithColour(title, url, user.Accent)
	} else {
		url := dctools.MemberAvatarURL(*member, guildID)
		url = dctools.ResizeImage(url, 4096)

		colour, _ := ctx.State.MemberColor(guildID, member.User.ID)
		if dctools.ColourInvalid(colour) {
			colour = discord.NullColor
		}

		embed = cmdutil.ImageInfoEmbedWithColour(title, url, colour)
	}

	ctx.RespondEmbed(*embed)
}
