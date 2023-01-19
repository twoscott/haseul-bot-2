package roles

import (
	"time"

	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
)

const maxRoleSelectionAge = time.Minute * 5

var (
	db             *database.DB
	selectionCache *roleCache
)

func Init(rt *router.Router) {
	db = database.GetInstance()
	selectionCache = newRoleCache(maxRoleSelectionAge)
	go selectionCache.ClearJob(time.Minute)

	rt.AddSelectListener(handleRoleSelect)
	rt.AddButtonListener(handleRoleButton)
	rt.AddMemberJoinHandler(handleNewMember)

	rt.AddCommand(rolePicker)
	rolePicker.AddSubCommandGroup(rolePickerRoles)
	rolePickerRoles.AddSubCommand(rolePickerRolesAdd)
	rolePickerRoles.AddSubCommand(rolePickerRolesRemove)
	rolePickerRoles.AddSubCommand(rolePickerRolesList)

	rolePicker.AddSubCommandGroup(rolePickerTiers)
	rolePickerTiers.AddSubCommand(rolePickerTiersAdd)
	rolePickerTiers.AddSubCommand(rolePickerTiersRemove)
	rolePickerTiers.AddSubCommand(rolePickerTiersSend)
	rolePickerTiers.AddSubCommand(rolePickerTiersList)

	rt.AddCommand(joinRoles)
	joinRoles.AddSubCommand(joinRolesAdd)
	joinRoles.AddSubCommand(joinRolesRemove)
	joinRoles.AddSubCommand(joinRolesClear)
	joinRoles.AddSubCommand(joinRolesList)
}
