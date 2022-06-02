package router

import (
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
	Options   discord.CommandInteractionOptions
	Deferred  bool
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

	var flags api.InteractionResponseFlags = 0
	if ctx.Ephemeral {
		flags |= api.EphemeralResponse
	}
	if ctx.Deferred {
		_, err := dctools.FollowupRespond(
			ctx.State, ctx.Interaction, api.InteractionResponseData{
				Content: option.NewNullableString(content),
				Embeds:  &embeds,
				Flags:   flags,
			},
		)
		return err
	}

	return dctools.MessageRespond(
		ctx.State, ctx.Interaction, api.InteractionResponseData{
			Content: option.NewNullableString(content),
			Embeds:  &embeds,
			Flags:   flags,
		},
	)
}

// RespondText responds to the supplied command with the supplied
// content.
func (ctx CommandCtx) RespondText(content string) error {
	var flags api.InteractionResponseFlags = 0
	if ctx.Ephemeral {
		flags |= api.EphemeralResponse
	}
	if ctx.Deferred {
		_, err := dctools.FollowupRespond(
			ctx.State, ctx.Interaction, api.InteractionResponseData{
				Content: option.NewNullableString(content),
				Flags:   flags,
			},
		)
		return err
	}

	return dctools.MessageRespond(
		ctx.State, ctx.Interaction, api.InteractionResponseData{
			Content: option.NewNullableString(content),
			Flags:   flags,
		},
	)
}

// RespondEmbed responds to the supplied command with the
// supplied embed(s).
func (ctx CommandCtx) RespondEmbed(embeds ...discord.Embed) error {
	var flags api.InteractionResponseFlags = 0
	if ctx.Ephemeral {
		flags |= api.EphemeralResponse
	}
	if ctx.Deferred {
		_, err := dctools.FollowupRespond(
			ctx.State, ctx.Interaction, api.InteractionResponseData{
				Embeds: &embeds,
				Flags:  flags,
			},
		)
		return err
	}

	return dctools.MessageRespond(
		ctx.State, ctx.Interaction, api.InteractionResponseData{
			Embeds: &embeds,
			Flags:  flags,
		},
	)
}

// RespondSuccess responds to a command with the provided content,
// prepended with a check emoji.
func (ctx CommandCtx) RespondSuccess(content string) error {
	var flags api.InteractionResponseFlags = 0
	succMsg := Success(content).String()
	if ctx.Ephemeral {
		flags |= api.EphemeralResponse
	}
	if ctx.Deferred {
		_, err := dctools.FollowupRespond(
			ctx.State, ctx.Interaction, api.InteractionResponseData{
				Content: option.NewNullableString(succMsg),
				Flags:   flags,
			},
		)
		return err
	}

	return dctools.MessageRespond(
		ctx.State, ctx.Interaction, api.InteractionResponseData{
			Content: option.NewNullableString(succMsg),
			Flags:   flags,
		},
	)
}

// RespondWarning responds to a command with the provided content,
// prepended with a warning emoji.
func (ctx CommandCtx) RespondWarning(content string) error {
	var flags api.InteractionResponseFlags = 0
	warnMsg := Warning(content).String()
	if ctx.Ephemeral {
		flags |= api.EphemeralResponse
	}
	if ctx.Deferred {
		_, err := dctools.FollowupRespond(
			ctx.State, ctx.Interaction, api.InteractionResponseData{
				Content: option.NewNullableString(warnMsg),
				Flags:   flags,
			},
		)
		return err
	}

	return dctools.MessageRespond(
		ctx.State, ctx.Interaction, api.InteractionResponseData{
			Content: option.NewNullableString(warnMsg),
			Flags:   flags,
		},
	)
}

// RespondError responds to a command with the provided content,
// prepended with a cross emoji.
func (ctx CommandCtx) RespondError(content string) error {
	var flags api.InteractionResponseFlags = 0
	errMsg := Error(content).String()
	if ctx.Ephemeral {
		flags |= api.EphemeralResponse
	}
	if ctx.Deferred {
		_, err := dctools.FollowupRespond(
			ctx.State, ctx.Interaction, api.InteractionResponseData{
				Content: option.NewNullableString(errMsg),
				Flags:   flags,
			},
		)
		return err
	}

	return dctools.MessageRespond(
		ctx.State, ctx.Interaction, api.InteractionResponseData{
			Content: option.NewNullableString(errMsg),
			Flags:   flags,
		},
	)
}

// RespondGenericError responds to a command with a
// generic error message.
func (ctx CommandCtx) RespondGenericError() error {
	var flags api.InteractionResponseFlags = 0
	errMsg := Error("Error occurred during command execution.").String()
	if ctx.Ephemeral {
		flags |= api.EphemeralResponse
	}
	if ctx.Deferred {
		_, err := dctools.FollowupRespond(
			ctx.State, ctx.Interaction, api.InteractionResponseData{
				Content: option.NewNullableString(errMsg),
				Flags:   flags,
			},
		)
		return err
	}

	return dctools.MessageRespond(
		ctx.State, ctx.Interaction, api.InteractionResponseData{
			Content: option.NewNullableString(errMsg),
			Flags:   flags,
		},
	)
}

// RespondConfirmationButtons responds to a command with a message, with
// both pager buttons and a confirmation button attached
func (ctx CommandCtx) RespondConfirmationButtons(
	content string, embeds ...discord.Embed) error {

	var flags api.InteractionResponseFlags = 0
	if ctx.Ephemeral {
		flags |= api.EphemeralResponse
	}
	if ctx.Deferred {
		return dctools.FollowupRespondConfirmationButtons(
			ctx.State, ctx.Interaction, api.InteractionResponseData{
				Content: option.NewNullableString(content),
				Embeds:  &embeds,
				Flags:   flags,
			},
		)
	}

	return dctools.MessageRespondConfirmationButtons(
		ctx.State, ctx.Interaction, api.InteractionResponseData{
			Content: option.NewNullableString(content),
			Embeds:  &embeds,
			Flags:   flags,
		},
	)
}

// RespondConfirmationButtons responds to a command with a message, with
// both pager buttons and a confirmation button attached
func (ctx CommandCtx) RespondPagerButtons(
	content string, embeds ...discord.Embed) error {

	var flags api.InteractionResponseFlags = 0
	if ctx.Ephemeral {
		flags |= api.EphemeralResponse
	}
	if ctx.Deferred {
		return dctools.FollowupRespondPagerButtons(
			ctx.State, ctx.Interaction, api.InteractionResponseData{
				Content: option.NewNullableString(content),
				Embeds:  &embeds,
				Flags:   flags,
			},
		)
	}

	return dctools.MessageRespondPagerButtons(
		ctx.State, ctx.Interaction, api.InteractionResponseData{
			Content: option.NewNullableString(content),
			Embeds:  &embeds,
			Flags:   flags,
		},
	)
}

// RespondCmdMessage responds with the pre-defined error, warning, or
// success command response.
func (ctx CommandCtx) RespondCmdMessage(response CmdResponse) error {
	return ctx.RespondText(response.String())
}

// RespondWithPaging responds to a slash command with a message pager
func (ctx CommandCtx) RespondPaging(messagePages []MessagePage) error {

	return ctx.respondWithPaging(messagePages, false)
}

// RespondConfirmationPaging responds to a slash command with a message pager
// and a confirmation button
func (ctx CommandCtx) RespondConfirmationPaging(
	messagePages []MessagePage) error {

	return ctx.respondWithPaging(messagePages, true)
}

func (ctx CommandCtx) respondWithPaging(
	messagePages []MessagePage,
	confirmation bool) error {

	if len(messagePages) < 2 {
		return ctx.RespondSimple(
			messagePages[0].Content,
			messagePages[0].Embeds...,
		)
	}

	var err error
	if confirmation {
		err = ctx.RespondConfirmationButtons(
			messagePages[0].Content,
			messagePages[0].Embeds...,
		)
	} else {
		err = ctx.RespondPagerButtons(
			messagePages[0].Content,
			messagePages[0].Embeds...,
		)
	}

	if err != nil {
		return err
	}

	return ctx.AddButtonPager(ctx.Interaction, messagePages)
}
