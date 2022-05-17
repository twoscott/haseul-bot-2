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
	// PagerButtons is a slice of discord components for use with button paging.
	PagerButtons = []discord.Component{
		discord.ButtonComponent{
			Label:    "First",
			CustomID: ButtonIDFirstPage,
			Style:    discord.SecondaryButton,
		},
		discord.ButtonComponent{
			Label:    "Prev",
			CustomID: ButtonIDPrevPage,
			Style:    discord.PrimaryButton,
		},
		discord.ButtonComponent{
			Label:    "Next",
			CustomID: ButtonIDNextPage,
			Style:    discord.PrimaryButton,
		},
		discord.ButtonComponent{
			Label:    "Last",
			CustomID: ButtonIDLastPage,
			Style:    discord.SecondaryButton,
		},
	}
	CheckButton = discord.ButtonComponent{
		Emoji: &discord.ButtonEmoji{
			ID:   ButtonCheckEmoji().ID,
			Name: ButtonCheckEmoji().Name,
		},
		CustomID: ButtonIDConfirm,
		Style:    discord.SuccessButton,
	}
)

// ReplyWithMessagePager sends a message reply with Button Pager buttons
// attached.
func ReplyWithMessagePager(
	st *state.State,
	msg *discord.Message,
	content string,
	embeds ...discord.Embed) (*discord.Message, error) {

	return replyWithPager(st, msg, content, false, embeds...)
}

// ReplyWithConfirmationPager sends a message reply with Button Pager buttons
// attached, and a confirmation button to confirm a page.
func ReplyWithConfirmationPager(
	st *state.State,
	msg *discord.Message,
	content string,
	embeds ...discord.Embed) (*discord.Message, error) {

	return replyWithPager(st, msg, content, true, embeds...)
}

func replyWithPager(
	st *state.State,
	msg *discord.Message,
	content string,
	confirm bool,
	embeds ...discord.Embed) (*discord.Message, error) {

	rowComponents := make([]discord.Component, len(PagerButtons))
	copy(rowComponents, PagerButtons)
	if confirm {
		rowComponents = append(rowComponents, CheckButton)
	}

	return ReplyNoPingComplex(st, msg, api.SendMessageData{
		Content: content,
		Embeds:  embeds,
		Components: []discord.Component{
			discord.ActionRowComponent{Components: rowComponents},
		},
	})
}

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

	rowComponents := make([]discord.Component, len(PagerButtons))
	copy(rowComponents, PagerButtons)
	if confirm {
		rowComponents = append(rowComponents, CheckButton)
	}

	return st.SendMessageComplex(channelID, api.SendMessageData{
		Content: content,
		Embeds:  embeds,
		Components: []discord.Component{
			discord.ActionRowComponent{Components: rowComponents},
		},
	})
}

//DisableAllButtons disables all buttons in a given message's components field.
func DisableAllButtons(st *state.State, msg *discord.Message) error {
	components := discord.UnwrapComponents(msg.Components)
	for _, c := range components {
		switch c := c.(type) {
		case *discord.ActionRowComponent:
			for _, c := range c.Components {
				switch c := c.(type) {
				case *discord.ButtonComponent:
					c.Disabled = true
				}
			}
		}
	}

	_, err := st.EditMessageComplex(msg.ChannelID, msg.ID, api.EditMessageData{
		Components: &components,
	})

	return err
}

func RemoveAllComponents(
	st *state.State,
	channelID discord.ChannelID,
	messageID discord.MessageID) error {

	_, err := st.EditMessageComplex(channelID, messageID,
		api.EditMessageData{
			Components: &[]discord.Component{},
		},
	)

	return err
}
