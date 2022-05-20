package router

import (
	"errors"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

type (
	// Router handles the routing of events to receiving functions.
	Router struct {
		State                  *state.State
		commands               CommandMap
		buttonPagers           ButtonPagerMap
		guildJoinListeners     []GuildJoinListener
		messageCreateListeners []MessageCreateListener
		mentionListeners       []MessageCreateListener
		startupListeners       []ReadyListener
	}

	CommandMap            map[string]*Command
	ButtonPagerMap        map[discord.MessageID]*ButtonPager
	MessageCreateListener func(*Router, *gateway.MessageCreateEvent)
	GuildJoinListener     func(*Router, *state.GuildJoinEvent)
	ReadyListener         func(*Router, *gateway.ReadyEvent)
)

// New returns a new instance of Router.
func New(state *state.State) *Router {
	return &Router{
		State:                  state,
		commands:               make(CommandMap),
		buttonPagers:           make(ButtonPagerMap),
		messageCreateListeners: make([]MessageCreateListener, 0),
		mentionListeners:       make([]MessageCreateListener, 0),
	}
}

// MustRegisterCommand adds a command object to the router.
func (rt *Router) MustRegisterCommand(cmd *Command) {
	registerCommandToDestination(rt.commands, cmd)
}

// HandleCommand handles an incoming message that could potentially match with
// a command.
func (rt *Router) HandleCommand(
	msg *gateway.MessageCreateEvent, args []string) {

	cmd, cmdLen := rt.findCommand(args)
	if cmd == nil || cmd.Run == nil {
		return
	}

	ctx := CommandCtx{
		Router: rt,
		Length: cmdLen,
		Member: msg.Member,
		Msg:    &msg.Message,
	}

	cmd.Execute(ctx, args[cmdLen:])
}

func (rt Router) findCommand(args []string) (*Command, int) {
	var depth = 0
	var result *Command
	var lastCmd *Command
	cmdsTracker := rt.commands

	for i, arg := range args {
		depth = i
		currCmd, ok := cmdsTracker[arg]
		if !ok {
			if lastCmd != nil {
				result = lastCmd
				depth = i
			}
			break
		}

		cmdsTracker = currCmd.SubCommands
		if cmdsTracker == nil || len(cmdsTracker) < 1 {
			result = currCmd
			depth = i + 1
			break
		}

		if i == len(args)-1 {
			result = currCmd
			depth = i + 1
			break
		}

		lastCmd = currCmd
	}

	return result, depth
}

// RegisterMessageHandler adds a function to receive all messages.
func (rt *Router) RegisterMessageHandler(
	messageCreateListener MessageCreateListener) {

	rt.messageCreateListeners = append(
		rt.messageCreateListeners, messageCreateListener,
	)
}

// HandleMessage routes a message create event to all listener functions
// registered to the router.
func (rt *Router) HandleMessage(msg *gateway.MessageCreateEvent) {
	for _, listener := range rt.messageCreateListeners {
		go listener(rt, msg)
	}
}

// RegisterGuildJoinHandler adds a function to receive all guild joins.
func (rt *Router) RegisterGuildJoinHandler(
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

// RegisterReadyListener adds a function to receive all ready events.
func (rt *Router) RegisterStartupListener(readyListener ReadyListener) {
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
func (rt *Router) AddButtonPager(options ButtonPagerOptions) error {
	if !options.MessageID.IsValid() {
		return errors.New(
			"No message ID was provided to add a button pager to",
		)
	}
	if len(options.Pages) < 2 {
		return nil
	}

	if _, ok := rt.buttonPagers[options.MessageID]; ok {
		return errors.New(
			"No more than one button pager can be assigned to a single message",
		)
	}

	buttonPager := newButtonPager(options)
	rt.buttonPagers[options.MessageID] = buttonPager

	go buttonPager.deleteAfterTimeout(rt)

	return nil
}

// HandleButton routes a button press to the relevant button pager.
func (rt *Router) HandleButtonPress(
	button *gateway.InteractionCreateEvent, data *discord.ButtonInteraction) {
	buttonPager, ok := rt.buttonPagers[button.Message.ID]
	if !ok {
		return
	}

	buttonPager.handleButton(rt, button, data)
}

func registerCommandToDestination(destination CommandMap, cmd *Command) {
	for _, c := range destination {
		if c == cmd {
			panic("Command cannot be registered twice")
		}
	}

	if cmd.Name == "" {
		panic("Failed to add a command with no name")
	}

	if destination == nil {
		panic("Command map must be initialised")
	}

	nameCheck, ok := destination[cmd.Name]
	if ok {
		log.Panicf("'%v is already registered to another command", nameCheck)
	}
	destination[cmd.Name] = cmd

	for _, alias := range cmd.Aliases {
		nameCheck, ok = destination[alias]
		if ok {
			log.Panicf(
				"'%v' is already registered to another command", nameCheck,
			)
		}
		destination[alias] = cmd
	}
}
