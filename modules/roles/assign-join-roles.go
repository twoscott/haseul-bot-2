package roles

import (
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

func handleNewMember(
	rt *router.Router, member discord.Member, guildID discord.GuildID) {

	roleIDs, err := db.Roles.GetAllGuildJoinRoles(guildID)
	if err != nil {
		log.Println(err)
		return
	}
	if len(roleIDs) < 1 {
		return
	}

	err = rt.State.ModifyMember(
		guildID,
		member.User.ID,
		api.ModifyMemberData{Roles: &roleIDs},
	)
	if err != nil {
		log.Println(err)
		return
	}
}
