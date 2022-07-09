package vlive

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var vliveRolesRemoveCommand = &router.SubCommand{
	Name:        "remove",
	Description: "Removes a mention role from a VLIVE feed",
	Handler: &router.CommandHandler{
		Executor:      vliveRoleRemoveExec,
		Autocompleter: vliveFeedRemoveCompleter,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:   "vlive-channel",
			Description:  "The VLIVE channel to search for VLIVE boards",
			Required:     true,
			Autocomplete: true,
		},
		&discord.IntegerOption{
			OptionName:   "vlive-board",
			Description:  "The VLIVE board to stop receiving VLIVE posts from",
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

func vliveRoleRemoveExec(ctx router.CommandCtx) {
	channelCode := ctx.Options.Find("vlive-channel").String()
	vChannel, cerr := fetchChannel(channelCode)
	if cerr != nil {
		ctx.RespondCmdMessage(cerr)
		return
	}

	boardID, err := ctx.Options.Find("vlive-board").IntValue()
	if err != nil {
		ctx.RespondWarning("Provided VLIVE board ID must be a number")
		return
	}
	board, cerr := fetchBoard(boardID)
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

	ok, err := db.VLIVE.RemoveMention(channel.ID, board.ID, role.ID)
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
				"%s is not set up to be mentioned in %s for %s VLIVE posts.",
				role.Mention(), channel.Mention(), vChannel.Name,
			),
		)
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf(
			"%s will no longer be mentioned in %s when there are new "+
				"%s VLIVE posts.",
			role.Mention(), channel.Mention(), vChannel.Name,
		),
	)
}
