package dctools

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

const (
	ButtonIDFirstPage = "FIRST_PAGE"
	ButtonIDPrevPage  = "PREV_PAGE"
	ButtonIDNextPage  = "NEXT_PAGE"
	ButtonIDLastPage  = "LAST_PAGE"
	ButtonIDConfirm   = "CONFIRM"
	ButtonIDTimeout   = "TIMEOUT"
)

var (
	// PagerActionRow is a set of buttons for use with button paging.
	PagerActionRow = discord.ActionRowComponent{
		&discord.ButtonComponent{
			Label:    "First",
			CustomID: ButtonIDFirstPage,
			Style:    discord.SecondaryButtonStyle(),
		},
		&discord.ButtonComponent{
			Label:    "Prev",
			CustomID: ButtonIDPrevPage,
			Style:    discord.PrimaryButtonStyle(),
		},
		&discord.ButtonComponent{
			Label:    "Next",
			CustomID: ButtonIDNextPage,
			Style:    discord.PrimaryButtonStyle(),
		},
		&discord.ButtonComponent{
			Label:    "Last",
			CustomID: ButtonIDLastPage,
			Style:    discord.SecondaryButtonStyle(),
		},
	}
	CheckButton = &discord.ButtonComponent{
		Label:    "Select",
		CustomID: ButtonIDConfirm,
		Style:    discord.SuccessButtonStyle(),
	}
)

// ReplyWithMessagePager sends a message reply with Button Pager buttons
// attached.
func SendWithMessagePager(
	st *state.State,
	channelID discord.ChannelID,
	content string,
	embeds ...discord.Embed) (*discord.Message, error) {

	return sendWithPager(st, channelID, content, false, embeds...)
}

// ReplyWithConfirmationPager sends a message reply with Button Pager buttons
// attached, and a confirmation button to confirm a page.
func SendWithConfirmationPager(
	st *state.State,
	channelID discord.ChannelID,
	content string,
	embeds ...discord.Embed) (*discord.Message, error) {

	return sendWithPager(st, channelID, content, true, embeds...)
}

func sendWithPager(
	st *state.State,
	channelID discord.ChannelID,
	content string,
	confirm bool,
	embeds ...discord.Embed) (*discord.Message, error) {

	rowComponents := make(discord.ActionRowComponent, len(PagerActionRow))
	copy(rowComponents, PagerActionRow)
	if confirm {
		rowComponents = append(rowComponents, CheckButton)
	}

	return st.SendMessageComplex(channelID, api.SendMessageData{
		Content:    content,
		Embeds:     embeds,
		Components: discord.Components(&rowComponents),
	})
}

// ReplyWithMessagePager sends a message reply with Button Pager buttons
// attached.
func ReplyWithMessagePager(
	st *state.State,
	msg *discord.Message,
	content string,
	embeds ...discord.Embed) (*discord.Message, error) {

	return replyWithPager(st, msg, false, content, embeds...)
}

// ReplyWithConfirmationPager sends a message reply with Button Pager buttons
// attached, and a confirmation button to confirm a page.
func ReplyWithConfirmationPager(
	st *state.State,
	msg *discord.Message,
	content string,
	embeds ...discord.Embed) (*discord.Message, error) {

	return replyWithPager(st, msg, true, content, embeds...)
}

func replyWithPager(
	st *state.State,
	msg *discord.Message,
	confirm bool,
	content string,
	embeds ...discord.Embed) (*discord.Message, error) {

	rowComponents := make(discord.ActionRowComponent, len(PagerActionRow))
	copy(rowComponents, PagerActionRow)
	if confirm {
		rowComponents = append(rowComponents, CheckButton)
	}

	return ReplyNoPingComplex(st, msg, api.SendMessageData{
		Content:    content,
		Embeds:     embeds,
		Components: discord.Components(&rowComponents),
	})
}

// MessageRespondPagerButtons responds to a slash command with a message, with
// pager buttons attached.
func MessageRespondPagerButtons(
	st *state.State,
	interaction *discord.InteractionEvent,
	data api.InteractionResponseData) error {

	return messageRespondPagerButtons(
		st, interaction, false, false, data,
	)
}

// MessageRespondPagerButtons responds to a slash command with a message, with
// both pager buttons and a confirmation button attached.
func MessageRespondConfirmationButtons(
	st *state.State,
	interaction *discord.InteractionEvent,
	data api.InteractionResponseData) error {

	return messageRespondPagerButtons(
		st, interaction, true, false, data,
	)
}

// FollowupRespondPagerButtons responds to a deferred slash command with
// a message, with pager buttons attached.
func FollowupRespondPagerButtons(
	st *state.State,
	interaction *discord.InteractionEvent,
	data api.InteractionResponseData) error {

	return messageRespondPagerButtons(
		st, interaction, false, true, data,
	)
}

// FollowupRespondPagerButtons responds to a deferred slash command with
// a message, with both pager buttons and a confirmation button attached.
func FollowupRespondConfirmationButtons(
	st *state.State,
	interaction *discord.InteractionEvent,
	data api.InteractionResponseData) error {

	return messageRespondPagerButtons(
		st, interaction, true, true, data,
	)
}

func messageRespondPagerButtons(
	st *state.State,
	interaction *discord.InteractionEvent,
	confirm bool,
	followup bool,
	data api.InteractionResponseData) error {

	rowComponents := make(discord.ActionRowComponent, len(PagerActionRow))
	copy(rowComponents, PagerActionRow)
	if confirm {
		rowComponents = append(rowComponents, CheckButton)
	}

	data.Components = discord.ComponentsPtr(&rowComponents)

	if followup {
		_, err := FollowupRespond(st, interaction, data)
		return err
	}

	return MessageRespond(st, interaction, data)
}

// DisabledButtons returns a copy of the provided buttons, all disabled.
func DisabledButtons(
	buttons discord.ActionRowComponent) discord.ActionRowComponent {

	newButtons := make(discord.ActionRowComponent, len(buttons))
	copy(newButtons, buttons)

	for i, c := range newButtons {
		switch b := c.(type) {
		case *discord.ButtonComponent:
			newButton := *b
			newButton.Disabled = true
			newButtons[i] = &newButton
		}
	}

	return newButtons
}
