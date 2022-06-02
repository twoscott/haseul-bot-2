package twitter

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var twtRolesRemoveCommand = &router.SubCommand{
	Name:        "remove",
	Description: "Removes a mention role from a Twitter feed",
	Handler: &router.CommandHandler{
		Executor:      twtRoleRemoveExec,
		Autocompleter: dbTwitterCompleter,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:   "twitter",
			Description:  "The Twitter user of the target feed",
			Required:     true,
			Autocomplete: true,
		},
		&discord.ChannelOption{
			OptionName:  "channel",
			Description: "The channel of the target feed",
			Required:    true,
			ChannelTypes: []discord.ChannelType{
				discord.GuildText,
				discord.GuildNews,
			},
		},
		&discord.RoleOption{
			OptionName:  "role",
			Description: "The role to stop mentioning when a Tweet is posted",
			Required:    true,
		},
	},
}

func twtRoleRemoveExec(ctx router.CommandCtx) {
	screenName := ctx.Options.Find("user").String()
	user, cerr := fetchUser(screenName)
	if cerr != nil {
		ctx.RespondCmdMessage(cerr)
		return
	}

	snowflake, _ := ctx.Options.Find("channel").SnowflakeValue()
	channelID := discord.ChannelID(snowflake)
	if !channelID.IsValid() {
		ctx.RespondWarning(
			"Malformed Discord channel provided.",
		)
		return
	}

	channel, err := ctx.State.Channel(channelID)
	if err != nil {
		ctx.RespondWarning("Invalid Discord channel provided.")
		return
	}

	snowflake, _ = ctx.Options.Find("role").SnowflakeValue()
	roleID := discord.RoleID(snowflake)

	role, err := ctx.State.Role(ctx.Interaction.GuildID, roleID)
	if err != nil {
		ctx.RespondWarning(
			"Invalid role ID provided.",
		)
		return
	}

	ok, err := db.Twitter.RemoveMention(channel.ID, user.ID, role.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			fmt.Sprintf(
				"Error occurred while removing %s from the database.",
				role.Mention(),
			),
		)
		return
	}
	if !ok {
		ctx.RespondWarning(
			fmt.Sprintf(
				"%s is not set up to be mentioned in %s.",
				role.Mention(), channel.Mention(),
			),
		)
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf(
			"%s will no longer be mentioned in %s when @%s tweets.",
			role.Mention(), channel.Mention(), user.ScreenName,
		),
	)
}
