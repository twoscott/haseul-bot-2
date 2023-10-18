package server

import (
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
)

var db *database.DB

func Init(rt *router.Router) {
	db = database.GetInstance()

	rt.AddGuildJoinHandler(onServerJoin)

	rt.AddCommand(serverCommand)
	serverCommand.AddSubCommand(serverBannerCommand)
	serverCommand.AddSubCommand(serverIconCommand)
	serverCommand.AddSubCommand(serverInfoCommand)
}

func onServerJoin(_ *router.Router, join *state.GuildJoinEvent) {
	db.Guilds.Add(join.Guild.ID)
}
