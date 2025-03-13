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

// TODO: clean up whole router, can probably have a single slice of listener
//       interfaces and check the types at runtime.

type (
	// Router handles the routing of events to receiving functions.
	Router struct {
		State                  *state.State
		commands               []*Command
		commandHandlers        CommandHandlers
		buttonPagers           ButtonPagerMap
		buttonListeners        []ButtonListener
		selectListeners        []SelectListener
		guildJoinListeners     []GuildJoinListener
		memberJoinListeners    []MemberJoinListener
		memberLeaveListeners   []MemberLeaveListener
		messageCreateListeners []MessageCreateListener
		messageDeleteListeners []MessageDeleteListener
		messageUpdateListeners []MessageUpdateListener
		mentionListeners       []MessageCreateListener
		startupListeners       []ReadyListener
	}

	CommandHandlers map[string]*CommandHandler
	ModalHandlers   map[discord.ComponentID]*CommandHandler
	ButtonPagerMap  map[discord.InteractionID]*ButtonPager
	SelectListener  func(
		*Router,
		*discord.InteractionEvent,
		*discord.StringSelectInteraction,
	)
	ButtonListener func(
		*Router,
		*discord.InteractionEvent,
		*discord.ButtonInteraction,
	)
	MessageCreateListener func(*Router, discord.Message, *discord.Member)
	MessageDeleteListener func(*Router, discord.Message)
	MessageUpdateListener func(
		*Router, discord.Message, discord.Message, *discord.Member,
	)
	GuildJoinListener   func(*Router, *state.GuildJoinEvent)
	MemberJoinListener  func(*Router, discord.Member, discord.GuildID)
	MemberLeaveListener func(*Router, discord.User, discord.GuildID)
	ReadyListener       func(*Router, *gateway.ReadyEvent)
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

	itx := InteractionCtx{
		Router:      rt,
		Interaction: interaction,
		Ephemeral:   handler.Ephemeral,
	}

	ctx := CommandCtx{
		InteractionCtx: &itx,
		Handler:        handler,
		Command:        command,
		Options:        userOptions,
	}

	if handler.Defer {
		itx.Defer()
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

	itx := InteractionCtx{
		Router:      rt,
		Interaction: interaction,
	}

	ctx := AutocompleteCtx{
		InteractionCtx: &itx,
		Options:        completionOptions,
		Focused:        *focusedOption,
	}

	handler.Autocomplete(ctx)
}

// GetRawCreateCommandData converts all commands stored in the router to
// Discord API create command data types.
func (rt Router) GetRawCreateCommandData() []api.CreateCommandData {
	newCommandData := make([]api.CreateCommandData, len(rt.commands))

	for i, cmd := range rt.commands {
		newCommandData[i] = *cmd.CreateData()
	}

	return newCommandData
}

// FindCommand finds a bot command by name.
func (rt Router) FindCommand(name string) *Command {
	for _, cmd := range rt.commands {
		if cmd.Name == name {
			return cmd
		}
	}

	return nil
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

	discordCommands, err := rt.State.BulkOverwriteCommands(app.ID, createData)
	for _, dcCmd := range discordCommands {
		cmd := rt.FindCommand(dcCmd.Name)
		if cmd == nil {
			log.Panic(
				"Unable to find command corresponding to a Discord command",
			)
		}

		cmd.discordID = dcCmd.ID
	}

	return err
}

// RegisterCommandHandlers maps command handler functions to their command
// triggers.
func (rt *Router) MustRegisterCommandHandlers() {
	for _, cmd := range rt.commands {
		for _, group := range cmd.SubCommandGroups {
			prefix := cmd.Name + " " + group.Name
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
		trigger := prefix + " " + cmd.Name
		rt.mustRegisterCommandHandler(trigger, cmd.Handler)
	}
}

func (rt *Router) mustRegisterCommandHandler(
	name string, handler *CommandHandler) {

	if handler == nil {
		log.Panicf("'%s' does not have a command handler", name)
	}

	nameCheck, ok := rt.commandHandlers[name]
	if ok {
		log.Panicf("'%v' is already registered to another command", nameCheck)
	}

	rt.commandHandlers[name] = handler
}

// AddMessageHandler adds a function to receive all messages.
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

// AddMessageDeleteHandler adds a function to receive all messages deleted.
func (rt *Router) AddMessageDeleteHandler(
	messageDeleteListener MessageDeleteListener) {

	rt.messageDeleteListeners = append(
		rt.messageDeleteListeners, messageDeleteListener,
	)
}

// HandleMessageDelete routes a message delete event to all listener functions
// registered to the router.
func (rt *Router) HandleMessageDelete(msg discord.Message) {

	for _, listener := range rt.messageDeleteListeners {
		go listener(rt, msg)
	}
}

// AddMessageDeleteHandler adds a function to receive all messages deleted.
func (rt *Router) AddMessageUpdateHandler(
	messageUpdateListener MessageUpdateListener) {

	rt.messageUpdateListeners = append(
		rt.messageUpdateListeners, messageUpdateListener,
	)
}

// HandleMessageDelete routes a message delete event to all listener functions
// registered to the router.
func (rt *Router) HandleMessageUpdate(
	old discord.Message, new discord.Message, member *discord.Member) {

	for _, listener := range rt.messageUpdateListeners {
		go listener(rt, old, new, member)
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

// AddButtonPager adds a button pager to the given message with the given pages.
func (rt *Router) AddButtonPager(
	interaction *discord.InteractionEvent, pages []MessagePage) error {

	if !interaction.ID.IsValid() {
		return errors.New(
			"no interaction ID was provided to add a button pager to",
		)
	}
	if len(pages) < 2 {
		return nil
	}

	if _, ok := rt.buttonPagers[interaction.ID]; ok {
		return errors.New(
			"no more than one button pager can be assigned to a single message",
		)
	}

	buttonPager := newButtonPager(interaction, pages)
	rt.buttonPagers[interaction.ID] = buttonPager

	go buttonPager.deleteAfterTimeout(rt)

	return nil
}

// AddButtonListener adds a function to receive all button press interactions.
func (rt *Router) AddButtonListener(buttonListener ButtonListener) {
	rt.buttonListeners = append(rt.buttonListeners, buttonListener)
}

// HandleButton routes a button press to the relevant button pager.
func (rt *Router) HandleButtonPress(
	interaction *discord.InteractionEvent, data *discord.ButtonInteraction) {

	for _, listener := range rt.buttonListeners {
		go listener(rt, interaction, data)
	}

	if interaction.Message.Interaction == nil {
		return
	}

	buttonPager, ok := rt.buttonPagers[interaction.Message.Interaction.ID]
	if !ok {
		return
	}

	buttonPager.handleButtonPress(rt, interaction, data)
}

// AddStartupListener adds a function to receive all select interactions.
func (rt *Router) AddSelectListener(selectListener SelectListener) {
	rt.selectListeners = append(rt.selectListeners, selectListener)
}

// HandleStartupEvent routes a select interaction to all listener functions
// registered to the router.
func (rt *Router) HandleSelect(
	interaction *discord.InteractionEvent,
	data *discord.StringSelectInteraction) {

	for _, listener := range rt.selectListeners {
		go listener(rt, interaction, data)
	}
}
