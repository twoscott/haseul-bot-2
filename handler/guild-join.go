package handler

import "github.com/diamondburned/arikawa/v3/state"

func (h *Handler) GuildJoin(guild *state.GuildJoinEvent) {
	h.HandleGuildJoin(guild)
}
