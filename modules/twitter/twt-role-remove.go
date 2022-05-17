package twitter

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var twtRoleRemoveCommand = &router.Command{
	Name:                "remove",
	Aliases:             []string{"rm", "delete", "del"},
	RequiredPermissions: discord.PermissionManageChannels | discord.PermissionManageRoles,
	IncludeAdmin:        true,
	UseTyping:           true,
	Run:                 twtRoleRemoveRun,
}

func twtRoleRemoveRun(ctx router.CommandCtx, args []string) {
	if len(args) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a Twitter user, Discord channel, "+
				"and a mention role to remove.",
		)
		return
	}
	if len(args) < 2 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a Discord channel and a mention role to remove.",
		)
		return
	}
	if len(args) < 3 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"Please provide a mention role to remove.",
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
				fmt.Sprintf("No role could be found with the name '%s'.", roleArg),
			)
			return
		}
	}

	role, err := ctx.State.Role(ctx.Msg.GuildID, roleID)
	if err != nil {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg, "Invalid role ID provided.")
		return
	}

	ok, err = db.Twitter.RemoveMention(channel.ID, user.ID, role.ID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			fmt.Sprintf(
				"Error occurred while removing %s from the database.",
				role.Mention(),
			),
		)
		return
	}
	if !ok {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			fmt.Sprintf(
				"%s is not set up to be mentioned in %s.",
				role.Mention(), channel.Mention(),
			),
		)
		return
	}

	dctools.ReplyWithSuccess(ctx.State, ctx.Msg,
		fmt.Sprintf(
			"%s will no longer be mentioned in %s when @%s tweets.",
			role.Mention(), channel.Mention(), user.ScreenName,
		),
	)
}
