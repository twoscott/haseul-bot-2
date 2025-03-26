package router

import (
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

// ButtonPager represents a listener for paged button presses on a message.
type ButtonPager struct {
	// Interaction is the interaction that invoked the pager.
	Interaction *discord.InteractionEvent
	// Pages consists of the pages that the attached buttons will
	// change between.
	Pages []MessagePage
	// PageNumber tracks the current page that the pager is on.
	PageNumber int
	// Timeout defines how long the pager will be active for. Once the defined
	// timeout period has elapsed, the buttons will be disabled and the pager
	// deleted from memory.
	Timeout time.Time
}

const listenTime = 5 * time.Minute

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
	// CheckButton is a button used for confirming the current page.
	CheckButton = &discord.ButtonComponent{
		Label:    "Select",
		CustomID: ButtonIDConfirm,
		Style:    discord.SuccessButtonStyle(),
	}
)

func PagerComponents() *discord.ContainerComponents {
	return getPagerComponents(false)
}

func ConfirmationComponents() *discord.ContainerComponents {
	return getPagerComponents(true)
}

func getPagerComponents(confirm bool) *discord.ContainerComponents {
	buttons := make(
		discord.ActionRowComponent, len(PagerActionRow), len(PagerActionRow)+1,
	)
	copy(buttons, PagerActionRow)

	if confirm {
		buttons = append(buttons, CheckButton)
	}

	return discord.ComponentsPtr(&buttons)
}

func newButtonPager(
	interaction *discord.InteractionEvent, pages []MessagePage) *ButtonPager {

	return &ButtonPager{
		Interaction: interaction,
		Pages:       pages,
		PageNumber:  0,
		Timeout:     time.Now().Add(listenTime),
	}
}

func (b ButtonPager) currentPage() *MessagePage {
	return &b.Pages[b.PageNumber]
}

func (b *ButtonPager) handleButtonPress(
	rt *Router,
	itx *discord.InteractionEvent,
	data *discord.ButtonInteraction) {

	user := itx.User
	if user == nil {
		user = &itx.Member.User
	}
	if user == nil {
		return
	}

	if user.ID != b.Interaction.SenderID() {
		// TODO: can itx be wrapped in an interaction-ctx before being passed,
		// so interaction response helper methods can be used instead of a
		// verbose state method, manually building an InteractionResponse object
		rt.State.RespondInteraction(itx.ID, itx.Token,
			api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: option.NewNullableString(
						Error(
							"You cannot interact with this message.",
						).String(),
					),
					Flags: discord.EphemeralMessage,
				},
			})
		return
	}

	if data.CustomID == ButtonIDConfirm {
		b.confirmPage(rt, itx)
		return
	}

	b.changePages(rt, itx, data)
}

func (b *ButtonPager) changePages(
	rt *Router,
	itx *discord.InteractionEvent,
	data *discord.ButtonInteraction) {

	startPage := b.PageNumber

	switch data.CustomID {
	case ButtonIDFirstPage:
		b.PageNumber = 0
	case ButtonIDLastPage:
		b.PageNumber = len(b.Pages) - 1
	case ButtonIDPrevPage:
		if b.PageNumber <= 0 {
			b.PageNumber = len(b.Pages) - 1
		} else {
			b.PageNumber--
		}
	case ButtonIDNextPage:
		if b.PageNumber >= len(b.Pages)-1 {
			b.PageNumber = 0
		} else {
			b.PageNumber++
		}
	}

	if b.PageNumber == startPage {
		rt.State.RespondInteraction(itx.ID, itx.Token,
			api.InteractionResponse{
				Type: api.DeferredMessageUpdate,
			},
		)
		return
	}

	newPage := b.currentPage().InteractionData()
	rt.State.RespondInteraction(itx.ID, itx.Token,
		api.InteractionResponse{
			Type: api.UpdateMessage,
			Data: newPage,
		},
	)
}

func (b ButtonPager) confirmPage(
	rt *Router, button *discord.InteractionEvent) {

	rt.State.RespondInteraction(button.ID, button.Token,
		api.InteractionResponse{
			Type: api.UpdateMessage,
			Data: &api.InteractionResponseData{
				Components: discord.ComponentsPtr(),
			},
		},
	)

	delete(rt.buttonPagers, b.Interaction.ID)
}

func (b ButtonPager) deleteAfterTimeout(rt *Router) {
	<-time.After(time.Until(b.Timeout))

	if _, ok := rt.buttonPagers[b.Interaction.ID]; !ok {
		return
	}

	disabledPagerButtons := dctools.DisabledButtons(PagerActionRow)

	rt.State.EditInteractionResponse(
		b.Interaction.AppID,
		b.Interaction.Token,
		api.EditInteractionResponseData{
			Components: discord.ComponentsPtr(&disabledPagerButtons),
		},
	)

	delete(rt.buttonPagers, b.Interaction.ID)
}

// MessagePage represents a page for button pagers.
type MessagePage struct {
	Content string
	Embeds  []discord.Embed
}

// InteractionData converts a message page to interaction response data that
// can be used to update an interaction response message.
func (p MessagePage) InteractionData() *api.InteractionResponseData {
	data := api.InteractionResponseData{
		Content:         option.NewNullableString(p.Content),
		AllowedMentions: dctools.NoMentions,
	}

	if len(p.Embeds) > 0 {
		data.Embeds = &p.Embeds
	}

	return &data
}
