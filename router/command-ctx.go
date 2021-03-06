package router

import (
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

// CommandCtx wraps router and includes data about the command interaction
// to be passed to the receiving command handler.
type CommandCtx struct {
	*Router
	Interaction *discord.InteractionEvent
	Command     *discord.CommandInteraction
	// Options contains options that were attached to the lowest level
	// command or sub command, this means it excludes sub command groups
	// or sub commands from the options.
	Options discord.CommandInteractionOptions
	// Deferred signals whether a deferred message response has been sent to
	// the interaction.
	Deferred bool
	// Ephemeral signals whether responses to the interacton should include
	// the ephemeral flag. This flag dictates that the response should only
	// be visiable to the initial interaction sender.
	Ephemeral bool
}

// Defer defers a command's response and if successful, sets the deferred
// state to true, making future command responses respond as followup
// messages instead of responses to message source.
func (ctx *CommandCtx) Defer() error {
	err := dctools.DeferResponse(ctx.State, ctx.Interaction)
	if err != nil {
		return err
	}

	ctx.Deferred = true
	return err
}

// Respond responds to the supplied command with the supplied
// response data.
func (ctx CommandCtx) Respond(data api.InteractionResponseData) error {
	if ctx.Ephemeral {
		data.Flags |= api.EphemeralResponse
	}
	if ctx.Deferred {
		_, err := dctools.FollowupRespond(ctx.State, ctx.Interaction, data)
		return err
	}

	return dctools.MessageRespond(ctx.State, ctx.Interaction, data)
}

// RespondSimple responds to the supplied command with the supplied
// content and embed(s).
func (ctx CommandCtx) RespondSimple(
	content string, embeds ...discord.Embed) error {

	return ctx.Respond(api.InteractionResponseData{
		Content: option.NewNullableString(content),
		Embeds:  &embeds,
	})
}

// RespondText responds to the supplied command with the supplied
// content.
func (ctx CommandCtx) RespondText(content string) error {
	return ctx.RespondSimple(content)
}

// RespondEmbed responds to the supplied command with the
// supplied embed(s).
func (ctx CommandCtx) RespondEmbed(embeds ...discord.Embed) error {
	return ctx.RespondSimple("", embeds...)
}

// RespondSuccess responds to a command with the provided content,
// prepended with a check emoji.
func (ctx CommandCtx) RespondSuccess(content string) error {
	return ctx.RespondText(Success(content).String())
}

// RespondWarning responds to a command with the provided content,
// prepended with a warning emoji.
func (ctx CommandCtx) RespondWarning(content string) error {
	return ctx.RespondText(Warning(content).String())
}

// RespondError responds to a command with the provided content,
// prepended with a cross emoji.
func (ctx CommandCtx) RespondError(content string) error {
	return ctx.RespondText(Error(content).String())
}

// RespondGenericError responds to a command with a
// generic error message.
func (ctx CommandCtx) RespondGenericError() error {
	errMsg := Error("Error occurred during command execution.").String()
	return ctx.RespondText(errMsg)
}

// RespondCmdMessage responds with the pre-defined error, warning, or
// success command response.
func (ctx CommandCtx) RespondCmdMessage(response CmdResponse) error {
	return ctx.RespondText(response.String())
}

// RespondConfirmationButtons responds to a command with a message, with
// both pager buttons and a confirmation button attached
func (ctx CommandCtx) RespondConfirmationButtons(
	content string, embeds ...discord.Embed) error {

	return ctx.Respond(api.InteractionResponseData{
		Content:    option.NewNullableString(content),
		Embeds:     &embeds,
		Components: ConfirmationComponents(),
	})
}

// RespondConfirmationButtons responds to a command with a message, with
// pager buttons attached.
func (ctx CommandCtx) RespondPagerButtons(
	content string, embeds ...discord.Embed) error {

	return ctx.Respond(api.InteractionResponseData{
		Content:    option.NewNullableString(content),
		Embeds:     &embeds,
		Components: PagerComponents(),
	})
}

// RespondWithPaging responds to a slash command with a message pager.
func (ctx CommandCtx) RespondPaging(messagePages []MessagePage) error {
	return ctx.respondWithPaging(messagePages, false)
}

// RespondConfirmationPaging responds to a slash command with a message pager
// and a confirmation button.
func (ctx CommandCtx) RespondConfirmationPaging(
	messagePages []MessagePage) error {

	return ctx.respondWithPaging(messagePages, true)
}

func (ctx CommandCtx) respondWithPaging(
	messagePages []MessagePage, confirmation bool) error {

	var (
		content = messagePages[0].Content
		embeds  = messagePages[0].Embeds
	)

	if len(messagePages) < 2 {
		return ctx.RespondSimple(content, embeds...)
	}

	var err error
	if confirmation {
		err = ctx.RespondConfirmationButtons(content, embeds...)
	} else {
		err = ctx.RespondPagerButtons(content, embeds...)
	}
	if err != nil {
		return err
	}

	return ctx.AddButtonPager(ctx.Interaction, messagePages)
}

func (ctx CommandCtx) ParseAccessibleChannel(
	channelID discord.ChannelID) (*discord.Channel, CmdResponse) {

	if !channelID.IsValid() {
		return nil, Warningf("Malformed Discord channel provided.")
	}

	channel, err := ctx.State.Channel(channelID)
	if dctools.ErrMissingAccess(err) {
		return nil, Warningf("I cannot access this channel.")
	}
	if err != nil {
		return nil, Warningf("Invalid Discord channel provided.")
	}
	if channel.GuildID != ctx.Interaction.GuildID {
		return nil, Warningf(
			"Channel provided must belong to this server.",
		)
	}
	if !dctools.IsTextChannel(channel.Type) {
		return nil, Warningf("Channel provided must be a text channel.")
	}

	return channel, nil
}

func (ctx CommandCtx) ParseSendableChannel(
	channelID discord.ChannelID) (*discord.Channel, CmdResponse) {

	channel, cerr := ctx.ParseAccessibleChannel(channelID)
	if cerr != nil {
		return channel, cerr
	}

	botUser, err := ctx.State.Me()
	if err != nil {
		log.Println(err)
		return nil, Errorf(
			"Error occurred checking my permissions in %s.",
			channel.Mention(),
		)
	}

	botPermissions, err := ctx.State.Permissions(channel.ID, botUser.ID)
	if err != nil {
		log.Println(err)
		return nil, Errorf(
			"Error occurred checking my permissions in %s.",
			channel.Mention(),
		)
	}

	neededPerms := dctools.PermissionsBitfield(
		discord.PermissionViewChannel,
		discord.PermissionSendMessages,
	)

	if !botPermissions.Has(neededPerms) {
		return nil, Errorf(
			"I do not have permission to send messages in %s!",
			channel.Mention(),
		)
	}

	return channel, nil
}
