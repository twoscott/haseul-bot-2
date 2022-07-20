package handler

import "github.com/diamondburned/arikawa/v3/gateway"

func (h *Handler) MemberLeave(ev *gateway.GuildMemberRemoveEvent) {
	h.Router.HandleMemberLeave(ev.User, ev.GuildID)
}
