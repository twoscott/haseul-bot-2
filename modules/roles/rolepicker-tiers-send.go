package roles

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var rolePickerTiersSend = &router.SubCommand{
	Name:        "send",
	Description: "Send a role picker tier to a channel",
	Handler: &router.CommandHandler{
		Executor:      rolePickerTiersSendExec,
		Autocompleter: tierNameCompleter,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:   "tier",
			Description:  "The tier to send to the channel for role picking",
			MaxLength:    option.NewInt(32),
			Required:     true,
			Autocomplete: true,
		},
		&discord.ChannelOption{
			OptionName:  "channel",
			Description: "The channel to send the role picker tier to",
			Required:    true,
			ChannelTypes: []discord.ChannelType{
				discord.GuildText,
				discord.GuildNews,
			},
		},
	},
}

func rolePickerTiersSendExec(ctx router.CommandCtx) {
	tierName := ctx.Options.Find("tier").String()

	tier, err := db.Roles.GetTierByName(ctx.Interaction.GuildID, tierName)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondError(
			fmt.Sprintf("The role tier '%s' does not exist.", tierName),
		)
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching role tiers.")
		return
	}

	channelVal, _ := ctx.Options.Find("channel").SnowflakeValue()
	channelID := discord.ChannelID(channelVal)
	if !channelID.IsValid() {
		ctx.RespondWarning(
			"Malformed Discord channel provided.",
		)
		return
	}

	channel, cerr := ctx.ParseSendableChannel(channelID)
	if cerr != nil {
		ctx.RespondCmdMessage(cerr)
		return
	}

	botUser, err := ctx.State.Me()
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while checking channel permissions.")
		return
	}

	botPermissions, err := ctx.State.Permissions(channel.ID, botUser.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while checking channel permissions.")
		return
	}

	if !botPermissions.Has(discord.PermissionManageRoles) {
		ctx.RespondWarning("I need permission to manage roles!")
		return
	}
	if !botPermissions.Has(discord.PermissionSendMessages) {
		ctx.RespondWarning(
			fmt.Sprintf(
				"I need permission to send messages to %s!",
				channel.ID,
			),
		)
		return
	}

	dbRoles, err := db.Roles.GetAllRolesByTier(tier.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching roles.")
		return
	}
	if len(dbRoles) < 1 {
		ctx.RespondWarning("This tier has no roles.")
		return
	}

	roles, err := ctx.State.Roles(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching role data.")
		return
	}

	options := make([]discord.SelectOption, 0, len(dbRoles))
	for _, dbRole := range dbRoles {
		for _, dcRole := range roles {
			if dcRole.ID == dbRole.ID {
				options = append(options, discord.SelectOption{
					Label:       dcRole.Name,
					Value:       dcRole.ID.String(),
					Description: dbRole.Description.String,
				})
			}
		}
	}

	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: util.TitleCase(tier.Name),
		},
		Description: tier.Description.String,
		Color:       dctools.EmbedBackColour,
	}

	_, err = ctx.State.SendMessageComplex(channel.ID, api.SendMessageData{
		Embeds: []discord.Embed{embed},
		Components: discord.Components(
			&discord.ActionRowComponent{
				&discord.StringSelectComponent{
					Options:     options,
					CustomID:    selectIDRoleSelect,
					Placeholder: "Select a role",
					ValueLimits: [2]int{1, len(options)},
				},
			},
			&discord.ActionRowComponent{
				&discord.ButtonComponent{
					Style:    discord.DangerButtonStyle(),
					CustomID: buttonIDRemoveSelectedRoles,
					Label:    "Remove Selected",
				},
				&discord.ButtonComponent{
					Style:    discord.SuccessButtonStyle(),
					CustomID: buttonIDAddSelectedRoles,
					Label:    "Add Selected",
				},
			},
		),
	})

	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while sending role tier picker.")
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf("Role tier sent to %s successfully", channel.Mention()),
	)
}
