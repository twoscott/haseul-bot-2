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
		Emoji: &discord.ComponentEmoji{
			ID:   ButtonCheckEmoji().ID,
			Name: ButtonCheckEmoji().Name,
		},
		CustomID: ButtonIDConfirm,
		Style:    discord.SuccessButtonStyle(),
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

// DisableAllButtons disables all buttons in a given message's components field.
func DisableAllButtons(st *state.State, msg *discord.Message) error {
	newComponents := make(discord.ContainerComponents, 0)
	copy(newComponents, msg.Components)
	for _, c := range msg.Components {
		switch c := c.(type) {
		case *discord.ActionRowComponent:
			for _, ic := range *c {
				switch ic := ic.(type) {
				case *discord.ButtonComponent:
					ic.Disabled = true
				}
			}

		}
	}

	_, err := st.EditMessageComplex(msg.ChannelID, msg.ID, api.EditMessageData{
		Components: &newComponents,
	})

	return err
}

// RemoveAllComponents removes all components from a message sent by the bot.
func RemoveAllComponents(
	st *state.State,
	channelID discord.ChannelID,
	messageID discord.MessageID) error {

	_, err := st.EditMessageComplex(channelID, messageID,
		api.EditMessageData{
			Components: &discord.ContainerComponents{},
		},
	)

	return err
}
