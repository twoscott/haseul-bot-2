package vlive

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var vliveRolesAddCommand = &router.SubCommand{
	Name:        "add",
	Description: "Adds a role to be mentioned when a VLIVE is posted to a feed",
	Handler: &router.CommandHandler{
		Executor:      vliveRoleAddExec,
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
			Description: "The role to mention when a Tweet is posted",
			Required:    true,
		},
	},
}

func vliveRoleAddExec(ctx router.CommandCtx) {
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

	roleIDs, err := db.VLIVE.GetMentionRoles(channel.ID, board.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while checking existing roles.",
		)
		return
	}
	if len(roleIDs) >= 10 {
		ctx.RespondWarning(
			"Only 10 roles can be assigned per VLIVE feed.",
		)
		return
	}

	canMention, err := checkMentionPermissions(ctx, channel, role)
	if err != nil {
		ctx.RespondError(
			fmt.Sprintf(
				"Error occurred checking my permissions to mention %s.",
				role.Mention(),
			),
		)
	}
	if !canMention {
		ctx.RespondWarning(
			fmt.Sprintf(
				"I do not have permission to mention %s!", role.Mention(),
			),
		)
	}

	_, err = db.VLIVE.GetFeed(channel.ID, board.ID)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondWarning(
			fmt.Sprintf(
				"There is no VLIVE feed currently set up for %s in %s.",
				vChannel.Name, channel.Mention(),
			),
		)
		return
	}

	ok, err := db.VLIVE.AddMention(channel.ID, board.ID, role.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			fmt.Sprintf(
				"Error occurred while adding %s to the database.",
				role.Mention(),
			),
		)
		return
	}
	if !ok {
		ctx.RespondWarning(
			fmt.Sprintf(
				"%s is already set up to be mentioned in %s for "+
					"%s VLIVE posts.",
				role.Mention(), channel.Mention(), vChannel.Name,
			),
		)
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf(
			"%s will now be mentioned in %s when there are new %s VLIVE posts.",
			role.Mention(), channel.Mention(), vChannel.Name,
		),
	)
}

func checkMentionPermissions(
	ctx router.CommandCtx,
	channel *discord.Channel,
	role *discord.Role) (bool, error) {

	if role.Mentionable {
		return true, nil
	}

	botUser, err := ctx.State.Me()
	if err != nil {
		return false, err
	}

	botPermissions, err := ctx.State.Permissions(channel.ID, botUser.ID)
	if err != nil {
		return false, err
	}

	neededPerms := discord.PermissionMentionEveryone
	if !botPermissions.Has(neededPerms) {
		return false, nil
	}

	return true, nil
}
