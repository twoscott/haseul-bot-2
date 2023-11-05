package roles

import (
	"log"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"golang.org/x/exp/slices"
)

type roleAction int

func (a roleAction) verb() string {
	switch a {
	case addRoleAction:
		return "add"
	case removeRoleAction:
		return "remove"
	default:
		return "modify"
	}
}

func (a roleAction) successTitle() string {
	switch a {
	case addRoleAction:
		return "Added"
	case removeRoleAction:
		return "Removed"
	default:
		return "Modified"
	}
}

func (a roleAction) failureTitle() string {
	switch a {
	case addRoleAction:
		return "Already Added"
	case removeRoleAction:
		return "Already Removed"
	default:
		return "Already Modified"
	}
}

const (
	addRoleAction roleAction = iota
	removeRoleAction
)

const (
	selectIDRoleSelect          = "ROLE_SELECT"
	buttonIDAddSelectedRoles    = "ADD_SELECTED_ROLES"
	buttonIDRemoveSelectedRoles = "REMOVE_SELECTED_ROLES"
)

func handleRoleSelect(
	rt *router.Router,
	interaction *discord.InteractionEvent,
	data *discord.StringSelectInteraction) {

	if data.CustomID != selectIDRoleSelect {
		return
	}

	roleIDs := make([]discord.RoleID, 0)
	for _, v := range data.Values {
		intID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Println(err)
			continue
		}

		roleIDs = append(roleIDs, discord.RoleID(intID))
	}

	selectionCache.SetSelection(interaction, roleIDs)

	rt.State.RespondInteraction(
		interaction.ID,
		interaction.Token,
		api.InteractionResponse{
			Type: api.MessageInteractionWithSource,
		},
	)
}

func handleRoleButton(
	rt *router.Router,
	interaction *discord.InteractionEvent,
	data *discord.ButtonInteraction) {

	var actionType roleAction
	switch data.CustomID {
	case buttonIDAddSelectedRoles:
		actionType = addRoleAction
	case buttonIDRemoveSelectedRoles:
		actionType = removeRoleAction
	default:
		return
	}

	targetRoleIDs := selectionCache.GetSelectedRoleIDs(interaction)
	if len(targetRoleIDs) < 1 {
		clearSelection(rt.State, interaction)
		dctools.FollowupRespond(rt.State, interaction,
			api.InteractionResponseData{
				Content: option.NewNullableString(
					router.Warningf(
						"Please select a role to %s.", actionType.verb(),
					).String(),
				),
				Flags: discord.EphemeralMessage,
			},
		)
		return
	}

	member, err := rt.State.Member(interaction.GuildID, interaction.SenderID())
	if err != nil {
		log.Println(err)
		dctools.MessageRespond(rt.State, interaction,
			api.InteractionResponseData{
				Content: option.NewNullableString(
					router.Error(
						"Error occurred while fetching roles.").String(),
				),
				Flags: discord.EphemeralMessage,
			},
		)
		return
	}

	ownedRolesIDs := member.RoleIDs
	successfulRoleIDs := make([]discord.RoleID, 0)
	failedRoleIDs := make([]discord.RoleID, 0)

	for _, targetID := range targetRoleIDs {
		foundRole := false
		roleIndex := -1
		for i, ownedID := range ownedRolesIDs {
			if targetID == ownedID {
				foundRole = true
				roleIndex = i
				break
			}
		}

		switch actionType {
		case addRoleAction:
			if foundRole {
				failedRoleIDs = append(failedRoleIDs, targetID)
			} else {
				ownedRolesIDs = append(ownedRolesIDs, targetID)
				successfulRoleIDs = append(successfulRoleIDs, targetID)
			}
		case removeRoleAction:
			if foundRole {
				ownedRolesIDs = slices.Delete(
					ownedRolesIDs, roleIndex, roleIndex+1,
				)
				successfulRoleIDs = append(successfulRoleIDs, targetID)
			} else {
				failedRoleIDs = append(failedRoleIDs, targetID)
			}
		}
	}

	err = rt.State.ModifyMember(
		interaction.GuildID,
		member.User.ID,
		api.ModifyMemberData{Roles: &ownedRolesIDs},
	)
	if dctools.ErrLackPermission(err) {
		log.Println(err)
		dctools.MessageRespondText(
			rt.State,
			interaction,
			router.Error(
				"I don't have permission to modify one or more of "+
					"the selected roles.").String(),
		)
		return
	}
	if err != nil {
		log.Println(err)
		dctools.MessageRespond(rt.State, interaction,
			api.InteractionResponseData{
				Content: option.NewNullableString(
					router.Error(
						"Error occurred while modifying roles.").String(),
				),
				Flags: discord.EphemeralMessage,
			},
		)
		return
	}

	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: "Modified Roles",
		},
		Color:  dctools.EmbedBackColour,
		Fields: []discord.EmbedField{},
	}

	successMentions := roleMentionsString(successfulRoleIDs)
	if successMentions != "" {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  actionType.successTitle(),
			Value: successMentions,
		})
	}
	failureMentions := roleMentionsString(failedRoleIDs)
	if failureMentions != "" {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  actionType.failureTitle(),
			Value: failureMentions,
		})
	}

	clearSelection(rt.State, interaction)
	rt.State.FollowUpInteraction(
		interaction.AppID,
		interaction.Token,
		api.InteractionResponseData{
			Embeds: &[]discord.Embed{embed},
			Flags:  discord.EphemeralMessage,
		},
	)
}

func clearSelection(
	st *state.State, interaction *discord.InteractionEvent) error {

	selectionCache.ClearSelection(interaction)

	return st.RespondInteraction(
		interaction.ID,
		interaction.Token,
		api.InteractionResponse{
			Type: api.UpdateMessage,
		},
	)
}

func roleMentionsString(roleIDs []discord.RoleID) (mentions string) {
	for _, id := range roleIDs {
		mentions += id.Mention() + " "
	}
	mentions = strings.TrimSpace(mentions)
	return
}
