package dctools

import (
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

	return
}

// MemberCanModifyRole returns whether a member has permission to modify a role.
func MemberCanModifyRole(
	st *state.State,
	guildID discord.GuildID,
	channelID discord.ChannelID,
	userID discord.UserID,
	roleID discord.RoleID) (can bool, err error) {

	g, err := st.Guild(guildID)
	if err != nil {
		return false, err
	}

	if IsOwner(*g, userID) {
		return true, nil
	}

	ch, err := st.Channel(channelID)
	if err != nil {
		return false, err
	}

	m, err := st.Member(guildID, userID)
	if err != nil {
		return false, err
	}

	r, err := st.Role(guildID, roleID)
	if err != nil {
		return false, err
	}

	p := discord.CalcOverwrites(*g, *ch, *m)
	if p.Has(discord.PermissionAdministrator) {
		return true, nil
	}

	if !p.Has(discord.PermissionManageRoles) {
		return false, nil
	}

	hr, err := MemberHighestRole(st, guildID, m.User.ID)
	if err != nil {
		return false, err
	}

	if r.Position > hr.Position {
		return
	}

	can = true
	return
}
