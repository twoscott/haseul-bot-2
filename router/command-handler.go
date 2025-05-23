package router

type CommandHandler struct {
	// Executor will be run when a chat input interaction is determined to be
	// aimed at the parent command.
	Executor func(CommandCtx)
	// Autocompleter will be run when an autocomplete interaction is determined
	// to be aimed at the parent command.
	Autocompleter func(AutocompleteCtx)
	// Defer defines whether the command will take longer than Discord's
	// pre-defined timeout of 3s to complete. If true, the command will
	// acknowledge the interaction first before calling the command's
	// handler.
	Defer bool
	// Ephemeral defines whether responses to the command will be ephemeral.
	// ephemeral messages are hidden from all but the user receiving
	// the response.
	Ephemeral bool
	
	// adminOnly determines whether the command should only be available in
	// the configured home guild, and only available for the configured bot
	// admin user. This value is set according to the top-most command or parent
	// command's AdminOnly field when the commands are initialised.
	adminOnly bool
}

// Execute runs the handler's Executor and handles any resulting panics.
func (h CommandHandler) Execute(ctx CommandCtx) {
	defer handleCommandPanic(ctx)
	h.Executor(ctx)
}

// Autocomplete runs the handler's Autocompleter and handles any
// resulting panics.
func (h CommandHandler) Autocomplete(ctx AutocompleteCtx) {
	defer handlePanic(ctx.State)
	h.Autocompleter(ctx)
}
