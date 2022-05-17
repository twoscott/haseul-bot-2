package router

import (
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

type (
	// ButtonPagerOptions is used to configure the creation of a button pager
	ButtonPagerOptions struct {
		AuthorID  discord.UserID
		ChannelID discord.ChannelID
		MessageID discord.MessageID
		Pages     []MessagePage
	}
	// MessagePage represents a page for button pagers.
	MessagePage struct {
		Content string
		Embeds  []discord.Embed
	}
)

// ButtonPager represents a listener for paged button presses on a message.
type ButtonPager struct {
	Timeout time.Time

	// AuthorID is the ID of the user who initiated the creation of the pager.
	AuthorID  discord.UserID
	ChannelID discord.ChannelID
	MessageID discord.MessageID

	Pages      []MessagePage
	PageNumber int
}

const listenTime = 3 * time.Minute

func newButtonPager(options ButtonPagerOptions) *ButtonPager {
	return &ButtonPager{
		Timeout:    time.Now().Add(listenTime),
		AuthorID:   options.AuthorID,
		ChannelID:  options.ChannelID,
		MessageID:  options.MessageID,
		Pages:      options.Pages,
		PageNumber: 0,
	}
}

func (b ButtonPager) currentPage() *MessagePage {
	return &b.Pages[b.PageNumber]
}

func (b *ButtonPager) handleButton(
	rt *Router, button *gateway.InteractionCreateEvent) {

	user := button.User
	if user == nil {
		user = &button.Member.User
	}
	if user == nil {
		return
	}

	if user.ID != b.AuthorID {
		return
	}

	if button.Data.CustomID == dctools.ButtonIDConfirm {
		b.confirmPage(rt, button)
		return
	}

	b.changePages(rt, button)
}

func (b *ButtonPager) changePages(
	rt *Router, button *gateway.InteractionCreateEvent) {

	startPage := b.PageNumber

	switch button.Data.CustomID {
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

	newPage := pageToInteractionData(b.currentPage())
	rt.State.RespondInteraction(button.ID, button.Token,
		api.InteractionResponse{
			Type: api.UpdateMessage,
			Data: newPage,
		},
	)
}

func (b *ButtonPager) confirmPage(
	rt *Router, button *gateway.InteractionCreateEvent) {

	rt.State.RespondInteraction(button.ID, button.Token,
		api.InteractionResponse{
			Type: api.UpdateMessage,
			Data: &api.InteractionResponseData{
				Components: &[]discord.Component{},
			},
		},
	)

	delete(rt.buttonPagers, b.MessageID)
}

func (b *ButtonPager) deleteAfterTimeout(rt *Router) {
	<-time.After(time.Until(b.Timeout))

	if _, ok := rt.buttonPagers[b.MessageID]; !ok {
		return
	}

	msg, err := rt.State.Message(b.ChannelID, b.MessageID)
	if err != nil {
		dctools.RemoveAllComponents(rt.State, b.ChannelID, b.MessageID)
	} else {
		dctools.DisableAllButtons(rt.State, msg)
	}

	delete(rt.buttonPagers, b.MessageID)
}

func pageToInteractionData(page *MessagePage) *api.InteractionResponseData {
	data := api.InteractionResponseData{
		Content:         option.NewNullableString(page.Content),
		AllowedMentions: dctools.NoMentions,
	}

	if len(page.Embeds) > 0 {
		data.Embeds = &page.Embeds
	}

	return &data
}
