package router

import (
	"errors"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

type (
	// Router handles the routing of events to receiving functions.
	Router struct {
		State                  *state.State
		commands               []*Command
		commandHandlers        CommandHandlers
		buttonPagers           ButtonPagerMap
		guildJoinListeners     []GuildJoinListener
		memberJoinListeners    []MemberJoinListener
		memberLeaveListeners   []MemberLeaveListener
		messageCreateListeners []MessageCreateListener
		mentionListeners       []MessageCreateListener
		startupListeners       []ReadyListener
	}

	CommandHandlers       map[string]*CommandHandler
	ButtonPagerMap        map[discord.InteractionID]*ButtonPager
	MessageCreateListener func(*Router, discord.Message, *discord.Member)
	GuildJoinListener     func(*Router, *state.GuildJoinEvent)
	MemberJoinListener    func(*Router, discord.Member, discord.GuildID)
	MemberLeaveListener   func(*Router, discord.User, discord.GuildID)
	ReadyListener         func(*Router, *gateway.ReadyEvent)
)

// New returns a new instance of Router.
func New(state *state.State) *Router {
	return &Router{
		State:                  state,
		commands:               make([]*Command, 0),
		commandHandlers:        make(CommandHandlers),
		buttonPagers:           make(ButtonPagerMap),
		memberJoinListeners:    make([]MemberJoinListener, 0),
		memberLeaveListeners:   make([]MemberLeaveListener, 0),
		messageCreateListeners: make([]MessageCreateListener, 0),
		mentionListeners:       make([]MessageCreateListener, 0),
	}
}

// TODO: add command counter; count how many times each command has been used.

// AddCommand adds a slash command to the router.
func (rt *Router) AddCommand(cmd *Command) {
	rt.commands = append(rt.commands, cmd)
}

// HandleCommand handles an incoming slash command event.
func (rt *Router) HandleCommand(
	interaction *discord.InteractionEvent,
	command *discord.CommandInteraction) {

	key := CommandInteractionKey(command)
	handler, ok := rt.commandHandlers[key]
	if !ok {
		err := errors.New("No command registered for '" + key + "'")
		log.Print(err)
		return
	}

	userOptions := dctools.CommandOptions(command)
	ctx := CommandCtx{
		Router:      rt,
		Interaction: interaction,
		Command:     command,
		Options:     userOptions,
		Ephemeral:   handler.Ephemeral,
	}

	if handler.Defer {
		ctx.Defer()
	}

	handler.Execute(ctx)
}

// HandleAutocomplete handles an autocomplete interaction.
func (rt *Router) HandleAutocomplete(
	interaction *discord.InteractionEvent,
	completion *discord.AutocompleteInteraction) {

	key := AutocompleteInteractionKey(completion)
	handler, ok := rt.commandHandlers[key]
	if !ok {
		err := errors.New("No command registered for '" + key + "'")
		log.Print(err)
		return
	}

	completionOptions := dctools.AutocompleteOptions(completion)
	focusedOption := dctools.FocusedOption(completion)
	ctx := AutocompleteCtx{
		Router:      rt,
		Interaction: interaction,
		Options:     completionOptions,
		Focused:     *focusedOption,
	}

	handler.Autocomplete(ctx)
}

// AddButtonPager adds a button pager to the given message with the given pages.
func (rt *Router) AddButtonPager(
	interaction *discord.InteractionEvent, pages []MessagePage) error {

	if !interaction.ID.IsValid() {
		return errors.New(
			"No interaction ID was provided to add a button pager to",
		)
	}
	if len(pages) < 2 {
		return nil
	}

	if _, ok := rt.buttonPagers[interaction.ID]; ok {
		return errors.New(
			"No more than one button pager can be assigned to a single message",
		)
	}

	buttonPager := newButtonPager(interaction, pages)
	rt.buttonPagers[interaction.ID] = buttonPager

	go buttonPager.deleteAfterTimeout(rt)

	return nil
}

// HandleButton routes a button press to the relevant button pager.
func (rt *Router) HandleButtonPress(
	button *discord.InteractionEvent, data *discord.ButtonInteraction) {

	buttonPager, ok := rt.buttonPagers[button.Message.Interaction.ID]
	if !ok {
		return
	}

	buttonPager.handleButtonPress(rt, button, data)
}

// GetRawCreateCommandData converts all commands stored in the router to
// Discord API create command data types.
func (rt *Router) GetRawCreateCommandData() []api.CreateCommandData {
	newCommandData := make([]api.CreateCommandData, len(rt.commands))

	for i, cmd := range rt.commands {
		newCommandData[i] = *cmd.CreateData()
	}

	return newCommandData
}

// AddCommandsToDiscord sends the defined commands to the Discord API,
// converted to their Discord API slash command types so that users can
// execute the commands.
func (rt *Router) AddCommandsToDiscord() error {
	createData := rt.GetRawCreateCommandData()

	app, err := rt.State.CurrentApplication()
	if err != nil {
		return err
	}

	_, err = rt.State.BulkOverwriteCommands(app.ID, createData)

	return err
}

// RegisterCommandHandlers maps command handler functions to their command
// triggers.
func (rt *Router) MustRegisterCommandHandlers() {
	for _, cmd := range rt.commands {
		for _, group := range cmd.SubCommandGroups {
			prefix := cmd.Name + "/" + group.Name
			rt.mustRegisterSubCommandHandlers(prefix, group.SubCommands)
		}

		prefix := cmd.Name
		rt.mustRegisterSubCommandHandlers(prefix, cmd.SubCommands)

		if len(cmd.SubCommandGroups) < 1 && len(cmd.SubCommands) < 1 {
			rt.mustRegisterCommandHandler(cmd.Name, cmd.Handler)
		}
	}
}

func (rt *Router) mustRegisterSubCommandHandlers(
	prefix string, subCommands []*SubCommand) {
	for _, cmd := range subCommands {
		trigger := prefix + "/" + cmd.Name
		rt.mustRegisterCommandHandler(trigger, cmd.Handler)
	}
}

func (rt *Router) mustRegisterCommandHandler(name string, handler *CommandHandler) {
	if handler == nil {
		log.Panicf("'%s' does not have a command handler", name)
	}

	nameCheck, ok := rt.commandHandlers[name]
	if ok {
		log.Panicf("'%v' is already registered to another command", nameCheck)
	}

	rt.commandHandlers[name] = handler
}

// RegisterMessageHandler adds a function to receive all messages.
func (rt *Router) AddMessageHandler(
	messageCreateListener MessageCreateListener) {

	rt.messageCreateListeners = append(
		rt.messageCreateListeners, messageCreateListener,
	)
}

// HandleMessage routes a message create event to all listener functions
// registered to the router.
func (rt *Router) HandleMessage(
	msg discord.Message, member *discord.Member) {

	for _, listener := range rt.messageCreateListeners {
		go listener(rt, msg, member)
	}
}

// RegisterGuildJoinHandler adds a function to receive all guild joins.
func (rt *Router) AddGuildJoinHandler(
	guildJoinListener GuildJoinListener) {

	rt.guildJoinListeners = append(rt.guildJoinListeners, guildJoinListener)
}

// HandleGuildJoin routes a guild join event to all listener functions
// registered to the router.
func (rt *Router) HandleGuildJoin(guild *state.GuildJoinEvent) {
	for _, listener := range rt.guildJoinListeners {
		go listener(rt, guild)
	}
}

// AddMemberJoinHandler adds a function to receive all member joins.
func (rt *Router) AddMemberJoinHandler(
	memberJoinListener MemberJoinListener) {

	rt.memberJoinListeners = append(rt.memberJoinListeners, memberJoinListener)
}

// HandleMemberJoin routers a member join event to all listener functions
// registered to the router.
func (rt *Router) HandleMemberJoin(
	member discord.Member, guildID discord.GuildID) {

	for _, listener := range rt.memberJoinListeners {
		go listener(rt, member, guildID)
	}
}

// AddMemberLeaveHandler adds a function to receive all member leaves.
func (rt *Router) AddMemberLeaveHandler(
	memberLeaveListener MemberLeaveListener) {

	rt.memberLeaveListeners = append(rt.memberLeaveListeners, memberLeaveListener)
}

// HandleMemberLeave routers a member leave event to all listener functions
// registered to the router.
func (rt *Router) HandleMemberLeave(
	user discord.User, guildID discord.GuildID) {

	for _, listener := range rt.memberLeaveListeners {
		go listener(rt, user, guildID)
	}
}

// AddStartupListener adds a function to receive all ready events.
func (rt *Router) AddStartupListener(readyListener ReadyListener) {
	rt.startupListeners = append(rt.startupListeners, readyListener)
}

// HandleStartupEvent routes a ready event to all listener functions
// registered to the router.
func (rt *Router) HandleStartupEvent(readyEvent *gateway.ReadyEvent) {
	for _, listener := range rt.startupListeners {
		go listener(rt, readyEvent)
	}
}
