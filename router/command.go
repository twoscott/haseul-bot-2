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
	"github.com/twoscott/haseul-bot-2/config"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

type (
	CommandHandler func(CommandCtx, []string)

	// Command represents a bot command.
	Command struct {
		Name    string
		Aliases []string

		RequiredPermissions discord.Permissions
		IncludeAdmin        bool
		UseTyping           bool

		Run         CommandHandler
		SubCommands map[string]*Command
	}

	// CommandCtx wraps router and is passed to command functions.
	CommandCtx struct {
		*Router
		Length int
		Member *discord.Member
		Msg    *discord.Message
	}
)

// Execute runs when a message's arguments match the command's.
func (c *Command) Execute(ctx CommandCtx, args []string) {
	if c.UseTyping {
		go ctx.State.Typing(ctx.Msg.ChannelID)
	}

	hasPerms, err := c.hasPermissions(ctx)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while attempting to check your permissions.",
		)
		return
	}
	if !hasPerms {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"You do not have permission to use this command!",
		)
		return
	}

	defer handleCommandPanic(ctx.State, ctx.Msg)
	c.Run(ctx, args)
}

// MustRegisterSubCommand adds a subcommand to a command.
func (c *Command) MustRegisterSubCommand(cmd *Command) {
	registerCommandToDestination(c.SubCommands, cmd)
}

func (c *Command) hasPermissions(ctx CommandCtx) (bool, error) {
	if c.RequiredPermissions == 0x0 {
		return true, nil
	}

	permissions, err := ctx.State.Permissions(
		ctx.Msg.ChannelID, ctx.Msg.Author.ID,
	)
	if err != nil {
		return false, err
	}

	var hasRequiredPerms bool
	if c.IncludeAdmin {
		hasRequiredPerms = dctools.HasAnyPermOrAdmin(
			permissions, c.RequiredPermissions,
		)
	} else {
		hasRequiredPerms = dctools.HasAnyPerm(
			permissions, c.RequiredPermissions,
		)
	}

	return hasRequiredPerms, nil
}

func handleCommandPanic(state *state.State, msg *discord.Message) {
	r := recover()
	if r == nil {
		return
	}

	errString := fmt.Sprintf("%v", r)

	dctools.ReplyWithError(state, msg, "Fatal Error Occurred.")
	log.Println("Recovered from command panic:", errString)
	debug.PrintStack()

	stackTrace := debug.Stack()
	logChannelID := config.GetInstance().Discord.LogChannelID
	if !logChannelID.IsValid() {
		return
	}

	stackReader := strings.NewReader(string(stackTrace))
	state.SendMessageComplex(logChannelID, api.SendMessageData{
		Content: dctools.Warning(errString),
		Files: []sendpart.File{
			{Name: "error.log", Reader: stackReader},
		},
	})
}
