package router

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

// CommandCtx wraps router and includes data about the command interaction
// to be passed to the receiving command handler.
type CommandCtx struct {
	*InteractionCtx
	Handler *CommandHandler
	Command *discord.CommandInteraction
	// Options contains options that were attached to the lowest level
	// command or sub command, this means it excludes sub command groups
	// or sub commands from the options.
	Options discord.CommandInteractionOptions
}

func (ctx CommandCtx) RespondWithModal(data api.InteractionResponseData) error {
	if data.CustomID != nil {
		customID := discord.ComponentID(data.CustomID.Val)
		ctx.modalHandlers[customID] = ctx.Handler
	}

	return dctools.ModalRespond(ctx.State, ctx.Interaction, data)
}
