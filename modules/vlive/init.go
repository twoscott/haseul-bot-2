package vlive

import (
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
	ptutil "github.com/twoscott/haseul-bot-2/utils/patreonutil"
)

var (
	db  *database.DB
	pat *ptutil.PatreonHelper
)

func Init(rt *router.Router) {
	db = database.GetInstance()
	pat = ptutil.GetPatreonHelper()

	rt.AddStartupListener(onStartup)

	rt.AddCommand(vliveCommand)
	vliveCommand.AddSubCommandGroup(vliveFeedsCommand)
	vliveFeedsCommand.AddSubCommand(vliveFeedsAddCommand)
	vliveFeedsCommand.AddSubCommand(vliveFeedsRemoveCommand)
	vliveFeedsCommand.AddSubCommand(vliveFeedsListCommand)
	vliveFeedsCommand.AddSubCommand(vliveFeedsClearCommand)

	vliveCommand.AddSubCommandGroup(vliveRolesCommand)
	vliveRolesCommand.AddSubCommand(vliveRolesAddCommand)
	vliveRolesCommand.AddSubCommand(vliveRolesRemoveCommand)
	vliveRolesCommand.AddSubCommand(vliveRolesListCommand)
	vliveRolesCommand.AddSubCommand(vliveRolesClearCommand)
}

func onStartup(rt *router.Router, _ *gateway.ReadyEvent) {
	startVLIVELoop(rt.State)
}
