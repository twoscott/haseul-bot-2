package roles

import (
	"fmt"
	"log"

	"github.com/twoscott/haseul-bot-2/router"
)

var joinRolesClear = &router.SubCommand{
	Name:        "clear",
	Description: "Clears all roles from being assigned to new members.",
	Handler: &router.CommandHandler{
		Executor: joinRolesClearExec,
	},
}

func joinRolesClearExec(ctx router.CommandCtx) {
	cleared, err := db.Roles.ClearGuildJoinRoles(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while clearing join roles.")
		return
	}
	if cleared == 0 {
		ctx.RespondWarning(
			"There are no join roles set up in this server to be cleared.",
		)
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf(
			"%d roles were cleared from being added to new members.",
			cleared,
		),
	)
}
