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
	button *discord.InteractionEvent,
	data *discord.ButtonInteraction) {

	user := button.User
	if user == nil {
		user = &button.Member.User
	}
	if user == nil {
		return
	}

	if user.ID != b.Interaction.SenderID() {
		return
	}

	if data.CustomID == dctools.ButtonIDConfirm {
		b.confirmPage(rt, button)
		return
	}

	b.changePages(rt, button, data)
}

func (b *ButtonPager) changePages(
	rt *Router,
	button *discord.InteractionEvent,
	data *discord.ButtonInteraction) {

	startPage := b.PageNumber

	switch data.CustomID {
	case dctools.ButtonIDFirstPage:
		b.PageNumber = 0
	case dctools.ButtonIDLastPage:
		b.PageNumber = len(b.Pages) - 1
	case dctools.ButtonIDPrevPage:
		if b.PageNumber <= 0 {
			b.PageNumber = len(b.Pages) - 1
		} else {
			b.PageNumber--
		}
	case dctools.ButtonIDNextPage:
		if b.PageNumber >= len(b.Pages)-1 {
			b.PageNumber = 0
		} else {
			b.PageNumber++
		}
	}

	if b.PageNumber == startPage {
		rt.State.RespondInteraction(button.ID, button.Token,
			api.InteractionResponse{
				Type: api.DeferredMessageUpdate,
			},
		)
		return
	}

	newPage := b.currentPage().InteractionData()
	rt.State.RespondInteraction(button.ID, button.Token,
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

	disabledPagerButtons := dctools.DisabledButtons(dctools.PagerActionRow)

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
