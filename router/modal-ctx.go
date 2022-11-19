package router

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

// ModalCtx wraps router and includes data about the modal submit
// interaction to be passed to the receiving modal handler.
type ModalCtx struct {
	*InteractionCtx
	Components discord.ContainerComponents
}
