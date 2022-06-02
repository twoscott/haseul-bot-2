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
)

type Command struct {
	Name                string
	Description         string
	Type                discord.CommandType
	Options             discord.CommandOptions
	RequiredPermissions *discord.Permissions
	SubCommandGroups    []*SubCommandGroup
	SubCommands         []*SubCommand
	Handler             *CommandHandler
}

// AddSubCommandGroup adds a group of sub commands to a slash command.
func (c *Command) AddSubCommandGroup(group *SubCommandGroup) {
	c.SubCommandGroups = append(c.SubCommandGroups, group)
}

// AddSubCommand adds a sub command to a slash command.
func (c *Command) AddSubCommand(command *SubCommand) {
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
}

// AddSubCommand adds a sub command to a command group.
func (g *SubCommandGroup) AddSubCommand(command *SubCommand) {
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
		return "/" + option.Name + commandString(&option.Options[0])
	case discord.SubcommandOptionType:
		return "/" + option.Name
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
	ctx.RespondError("Fatal Error Occurred during command execution.")
	log.Println("Recovered from command panic:", errString)
	debug.PrintStack()

	logPanicStack(ctx.State, errString)
}

func handleAutocompletePanic(ctx AutocompleteCtx) {
	r := recover()
	if r == nil {
		return
	}

	errString := fmt.Errorf("%v", r).Error()
	log.Println("Recovered from autocomplete panic:", errString)
	debug.PrintStack()

	logPanicStack(ctx.State, errString)
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
