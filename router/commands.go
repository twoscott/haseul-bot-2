package router

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/twoscott/haseul-bot-2/utils/botutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

type ParentCommand interface {
	ID() discord.CommandID
	NameReference() []string
}

type Command struct {
	Name                string
	Description         string
	Type                discord.CommandType
	Options             discord.CommandOptions
	RequiredPermissions *discord.Permissions
	SubCommandGroups    []*SubCommandGroup
	SubCommands         []*SubCommand
	Handler             *CommandHandler
	discordID           *discord.CommandID
}

// ID implements ParentCommand and facilitates fetching command mentions.
func (c Command) ID() discord.CommandID {
	if c.discordID == nil {
		return discord.NullCommandID
	}

	return *c.discordID
}

// NameReference implements ParentCommand and facilitiates
// fetching command mentions.
func (c Command) NameReference() []string {
	return []string{c.Name}
}

// Mention returns the mention string for a Discord message.
func (c Command) Mention() string {
	return dctools.CommandMention(c.ID(), c.NameReference()...)
}

// AddSubCommandGroup adds a group of sub commands to a slash command.
func (c *Command) AddSubCommandGroup(group *SubCommandGroup) {
	group.Parent = c
	c.SubCommandGroups = append(c.SubCommandGroups, group)
}

// AddSubCommand adds a sub command to a slash command.
func (c *Command) AddSubCommand(command *SubCommand) {
	command.Parent = c
	c.SubCommands = append(c.SubCommands, command)
}

// CreateData converts a command into its underlying create command data
// discord API type.
func (c Command) CreateData() *api.CreateCommandData {
	createData := api.CreateCommandData{
		Name:                     c.Name,
		Description:              c.Description,
		Type:                     c.Type,
		DefaultMemberPermissions: c.RequiredPermissions,
	}

	for _, group := range c.SubCommandGroups {
		createData.Options = append(createData.Options, group.CreateData())
	}

	for _, cmd := range c.SubCommands {
		createData.Options = append(createData.Options, cmd.CreateData())
	}

	createData.Options = append(createData.Options, c.Options...)

	return &createData
}

type SubCommandGroup struct {
	Name        string
	Description string
	SubCommands []*SubCommand
	Parent      *Command
}

// ID implements ParentCommand and facilitates fetching command mentions.
func (c SubCommandGroup) ID() discord.CommandID {
	if c.Parent == nil {
		return discord.NullCommandID
	}

	return c.Parent.ID()
}

// NameReference implements ParentCommand and facilitiates
// fetching command mentions.
func (c SubCommandGroup) NameReference() []string {
	return append(c.Parent.NameReference(), c.Name)
}

// Mention returns the mention string for a Discord message.
func (c SubCommandGroup) Mention() string {
	return dctools.CommandMention(c.ID(), c.NameReference()...)
}

// AddSubCommand adds a sub command to a command group.
func (g *SubCommandGroup) AddSubCommand(command *SubCommand) {
	command.Parent = g
	g.SubCommands = append(g.SubCommands, command)
}

// ToCreateData converts a sub command group into its underlying sub command
// group option API type.
func (g SubCommandGroup) CreateData() *discord.SubcommandGroupOption {
	optionData := discord.SubcommandGroupOption{
		OptionName:  g.Name,
		Description: g.Description,
	}

	for _, cmd := range g.SubCommands {
		optionData.Subcommands = append(
			optionData.Subcommands, cmd.CreateData(),
		)
	}

	return &optionData
}

type SubCommand struct {
	Name        string
	Description string
	// Defer defines whether the command will take longer than Discord's
	// pre-defined timeout of 3s to complete. If true, the command will
	// acknowledge the interaction first before calling the command's
	// handler.
	Options []discord.CommandOptionValue
	Handler *CommandHandler
	Parent  ParentCommand
}

func (c SubCommand) ID() discord.CommandID {
	if c.Parent == nil {
		return discord.NullCommandID
	}

	return c.Parent.ID()
}

// NameReference implements ParentCommand and facilitiates
// fetching command mentions.
func (c *SubCommand) NameReference() []string {
	return append(c.Parent.NameReference(), c.Name)
}

// Mention returns the mention string for a Discord message.
func (c SubCommand) Mention() string {
	return dctools.CommandMention(c.ID(), c.NameReference()...)
}

// ToCreateData converts a sub command nto its underlying sub command option
// API type.
func (c SubCommand) CreateData() *discord.SubcommandOption {
	return &discord.SubcommandOption{
		OptionName:  c.Name,
		Description: c.Description,
		Options:     c.Options,
	}
}

// CommandInteractionKey returns the string representing the command and
// subcommands of a command interaction as a single string,
// for use in the commandHandlers hash map.
func CommandInteractionKey(command *discord.CommandInteraction) string {
	if len(command.Options) < 1 {
		return command.Name
	}

	return command.Name + commandString(&command.Options[0])
}

func commandString(option *discord.CommandInteractionOption) string {
	switch option.Type {
	case discord.SubcommandGroupOptionType:
		return " " + option.Name + commandString(&option.Options[0])
	case discord.SubcommandOptionType:
		return " " + option.Name
	default:
		return ""
	}
}

func handleCommandPanic(ctx CommandCtx) {
	r := recover()
	if r == nil {
		return
	}

	errString := fmt.Errorf("%v", r).Error()
	ctx.RespondError("Fatal error occurred during command execution.")
	log.Println("Recovered from command panic:", errString)
	debug.PrintStack()

	logPanicStack(ctx.State, errString)
}

func handlePanic(st *state.State) {
	r := recover()
	if r == nil {
		return
	}

	errString := fmt.Errorf("%v", r).Error()
	log.Println("Recovered from autocomplete panic:", errString)
	debug.PrintStack()

	logPanicStack(st, errString)
}

func logPanicStack(st *state.State, errString string) {
	stackTrace := debug.Stack()
	stackReader := strings.NewReader(string(stackTrace))

	logData := api.SendMessageData{
		Content: Warning(errString).String(),
		Files: []sendpart.File{
			{Name: "error.log", Reader: stackReader},
		},
	}

	botutil.Log(st, logData)
}
