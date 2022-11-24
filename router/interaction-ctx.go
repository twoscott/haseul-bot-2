package router

import (
	"io"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

// InteractionCtx wraps router and includes data about the interaction
// to be passed to the receiving handler.
type InteractionCtx struct {
	*Router
	Interaction *discord.InteractionEvent
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
func (ctx *InteractionCtx) Defer() error {
	err := dctools.DeferResponse(ctx.State, ctx.Interaction)
	if err != nil {
		return err
	}

	ctx.Deferred = true
	return err
}

// Respond responds to the supplied command with the supplied
// response data.
func (ctx InteractionCtx) Respond(data api.InteractionResponseData) error {
	if ctx.Ephemeral {
		data.Flags |= discord.EphemeralMessage
	}
	if ctx.Deferred {
		_, err := dctools.FollowupRespond(ctx.State, ctx.Interaction, data)
		return err
	}

	return dctools.MessageRespond(ctx.State, ctx.Interaction, data)
}

// RespondSimple responds to the supplied command with the supplied
// content and embed(s).
func (ctx InteractionCtx) RespondSimple(
	content string, embeds ...discord.Embed) error {

	return ctx.Respond(api.InteractionResponseData{
		Content: option.NewNullableString(content),
		Embeds:  &embeds,
	})
}

// RespondText responds to the supplied command with the supplied content.
func (ctx InteractionCtx) RespondText(content string) error {
	return ctx.RespondSimple(content)
}

// RespondEmbed responds to the supplied command with the supplied embed(s).
func (ctx InteractionCtx) RespondEmbed(embeds ...discord.Embed) error {
	return ctx.RespondSimple("", embeds...)
}

// RespondFiles responds to the supplied command with the supplied file(s).
func (ctx InteractionCtx) RespondFiles(files []sendpart.File) error {
	return ctx.Respond(api.InteractionResponseData{Files: files})
}

// RespondFile responds to the supplied command with the supplied file.
func (ctx InteractionCtx) RespondFile(fileName string, data io.Reader) error {
	return ctx.RespondFiles([]sendpart.File{
		{Name: fileName, Reader: data},
	})
}

// RespondSuccess responds to a command with the provided content,
// prepended with a check emoji.
func (ctx InteractionCtx) RespondSuccess(content string) error {
	return ctx.RespondText(Success(content).String())
}

// RespondWarning responds to a command with the provided content,
// prepended with a warning emoji.
func (ctx InteractionCtx) RespondWarning(content string) error {
	return ctx.RespondText(Warning(content).String())
}

// RespondError responds to a command with the provided content,
// prepended with a cross emoji.
func (ctx InteractionCtx) RespondError(content string) error {
	return ctx.RespondText(Error(content).String())
}

// RespondGenericError responds to a command with a
// generic error message.
func (ctx InteractionCtx) RespondGenericError() error {
	errMsg := Error("Error occurred during command execution.").String()
	return ctx.RespondText(errMsg)
}

// RespondCmdMessage responds with the pre-defined error, warning, or
// success command response.
func (ctx InteractionCtx) RespondCmdMessage(response CmdResponse) error {
	return ctx.RespondText(response.String())
}

// RespondConfirmationButtons responds to a command with a message, with
// both pager buttons and a confirmation button attached
func (ctx InteractionCtx) RespondConfirmationButtons(
	content string, embeds ...discord.Embed) error {

	return ctx.Respond(api.InteractionResponseData{
		Content:    option.NewNullableString(content),
		Embeds:     &embeds,
		Components: ConfirmationComponents(),
	})
}

// RespondConfirmationButtons responds to a command with a message, with
// pager buttons attached.
func (ctx InteractionCtx) RespondPagerButtons(
	content string, embeds ...discord.Embed) error {

	return ctx.Respond(api.InteractionResponseData{
		Content:    option.NewNullableString(content),
		Embeds:     &embeds,
		Components: PagerComponents(),
	})
}

// RespondWithPaging responds to a slash command with a message pager.
func (ctx InteractionCtx) RespondPaging(messagePages []MessagePage) error {
	return ctx.respondWithPaging(messagePages, false)
}

// RespondConfirmationPaging responds to a slash command with a message pager
// and a confirmation button.
func (ctx InteractionCtx) RespondConfirmationPaging(
	messagePages []MessagePage) error {

	return ctx.respondWithPaging(messagePages, true)
}

func (ctx InteractionCtx) respondWithPaging(
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

func (ctx InteractionCtx) ParseAccessibleChannel(
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

func (ctx InteractionCtx) ParseSendableChannel(
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
