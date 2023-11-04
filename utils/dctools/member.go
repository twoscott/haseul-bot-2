package dctools

import (
	"errors"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

// MemberHighestRole returns the highest-positioned role that a member has.
func MemberHighestRole(
	st *state.State,
	guildID discord.GuildID,
	userID discord.UserID) (role *discord.Role, err error) {

	m, err := st.Member(guildID, userID)
	if err != nil {
		return
	}

	for _, id := range m.RoleIDs {
		r, err := st.Role(guildID, id)
		if err != nil {
			continue
		}

		if role == nil || r.Position > role.Position {
			role = r
		}
	}

	if role == nil {
		return nil, errors.New("no role found for user")
	}

	return
}

// MemberCanModifyRole returns whether a member has permission to modify a role.
func MemberCanModifyRole(
	st *state.State,
	guildID discord.GuildID,
	channelID discord.ChannelID,
	userID discord.UserID,
	roleID discord.RoleID) (bool, error) {

	p, err := st.Permissions(channelID, userID)
	if err != nil {
		return false, err
	}

	if !HasAnyPermOrAdmin(p, discord.PermissionManageRoles) {
		return false, err
	}

	hr, err := MemberHighestRole(st, guildID, userID)
	if err != nil {
		return false, err
	}

	r, err := st.Role(guildID, roleID)
	if err != nil {
		return false, err
	}

	if hr.Position < r.Position {
		return false, nil
	}

	return true, nil
}
