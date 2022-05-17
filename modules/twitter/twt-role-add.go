package twitter

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var twtRoleAddCommand = &router.Command{
	Name: "add",
	RequiredPermissions: 0 |
		discord.PermissionManageChannels |
		discord.PermissionManageRoles,
	IncludeAdmin: true,
	UseTyping:    true,
	Run:          twtRoleAddRun,
}

func twtRoleAddRun(ctx router.CommandCtx, args []string) {
	if len(args) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a Twitter user, Discord channel, "+
				"and a mention role to add.",
		)
		return
	}
	if len(args) < 2 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a Discord channel and a mention role to add.",
		)
		return
	}
	if len(args) < 3 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a mention role to add.",
		)
		return
	}

	screenName := parseUserIfURL(args[0])
	user, ok := fetchUser(ctx, screenName)
	if !ok {
		return
	}

	channel, ok := parseChannelArg(ctx, args[1])
	if !ok {
		return
	}

	roleIDs, err := db.Twitter.GetMentionRoles(channel.ID, user.ID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while checking existing roles.",
		)
		return
	}
	if len(roleIDs) >= 10 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Only 10 roles can be assigned per Twitter feed.",
		)
		return
	}

	roleArg := util.TrimArgs(ctx.Msg.Content, ctx.Length+2)
	roleID := dctools.ParseRoleID(roleArg)
	if !roleID.IsValid() {
		roles, err := ctx.State.Roles(ctx.Msg.GuildID)
		if err != nil {
			log.Println(err)
			dctools.ReplyWithError(ctx.State, ctx.Msg,
				"Error occurred while trying to find provided role name.",
			)
			return
		}

		roleID = dctools.FindRoleIDByName(roles, roleArg)
		if !roleID.IsValid() {
			dctools.ReplyWithWarning(ctx.State, ctx.Msg,
				fmt.Sprintf(
					"No role could be found with the name '%s'.", roleArg,
				),
			)
			return
		}
	}

	role, err := ctx.State.Role(ctx.Msg.GuildID, roleID)
	if err != nil {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Invalid role ID provided.",
		)
		return
	}

	canMention, err := checkMentionPermissions(ctx, channel, role)
	if err != nil {
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			fmt.Sprintf(
				"Error occurred checking my permissions to mention %s.",
				role.Mention(),
			),
		)
	}
	if !canMention {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			fmt.Sprintf(
				"I do not have permission to mention %s!", role.Mention(),
			),
		)
	}

	_, err = db.Twitter.GetFeed(channel.ID, user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			fmt.Sprintf(
				"There is no Twitter feed currently set up for @%s in %s.",
				user.ScreenName, channel.Mention(),
			),
		)
		return
	}

	ok, err = db.Twitter.AddMention(channel.ID, user.ID, role.ID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			fmt.Sprintf(
				"Error occurred while adding %s to the database.",
				role.Mention(),
			),
		)
		return
	}
	if !ok {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			fmt.Sprintf(
				"%s is already set up to be mentioned in %s.",
				role.Mention(), channel.Mention(),
			),
		)
		return
	}

	dctools.ReplyWithSuccess(ctx.State, ctx.Msg,
		fmt.Sprintf(
			"%s will now be mentioned in %s when @%s tweets.",
			role.Mention(), channel.Mention(), user.ScreenName,
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
