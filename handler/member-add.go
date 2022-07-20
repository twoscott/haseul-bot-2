package handler

import "github.com/diamondburned/arikawa/v3/gateway"

func (h *Handler) MemberJoin(ev *gateway.GuildMemberAddEvent) {
	h.Router.HandleMemberJoin(ev.Member, ev.GuildID)
}
