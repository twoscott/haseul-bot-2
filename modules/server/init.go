package server

import (
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
)

var db *database.DB

func Init(rt *router.Router) {
	db = database.GetInstance()

	rt.AddStartupListener(onStartup)
	rt.AddGuildJoinHandler(onServerJoin)

	rt.AddCommand(serverCommand)
	serverCommand.AddSubCommand(serverBannerCommand)
	serverCommand.AddSubCommand(serverIconCommand)
	serverCommand.AddSubCommand(serverInfoCommand)
}

func onServerJoin(_ *router.Router, join *state.GuildJoinEvent) {
	db.Guilds.Add(join.Guild.ID)
}

func onStartup(rt *router.Router, _ *gateway.ReadyEvent) {
	guilds, _ := rt.State.Guilds()
	for _, guild := range guilds {
		db.Guilds.Add(guild.ID)
	}
}
