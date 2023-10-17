package user

import (
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
)

var db *database.DB

func Init(rt *router.Router) {
	db = database.GetInstance()

	rt.AddMessageHandler(addXP)

	rt.AddCommand(levelsCommand)
	levelsCommand.AddSubCommand(levelsLeaderboardCommand)

	rt.AddCommand(repCommand)
	repCommand.AddSubCommand(repGiveCommand)
	repCommand.AddSubCommand(repStatusCommand)
	repCommand.AddSubCommand(repLeaderboardCommand)

	repCommand.AddSubCommandGroup(repStreaksCommand)
	repStreaksCommand.AddSubCommand(repStreaksLeaderboardCommand)

	rt.AddCommand(userCommand)
	userCommand.AddSubCommand(userAvatarCommand)
	userCommand.AddSubCommand(userBannerCommand)
	userCommand.AddSubCommand(userInfoCommand)
}
