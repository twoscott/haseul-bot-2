// Package dctools provides helper functions and constants relevant to Discord
// and its API.
package dctools

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

// NoMentions is a zero value AllowedMentions object which allows no mentions.
var NoMentions = new(api.AllowedMentions)

// ReplyNoPingComplex replies to a message with the provided message data,
// with mentions disabled.
func ReplyNoPingComplex(
	st *state.State,
	msg *discord.Message,
	data api.SendMessageData) (*discord.Message, error) {

	data.AllowedMentions = NoMentions
	data.Reference = &discord.MessageReference{MessageID: msg.ID}
	return st.SendMessageComplex(msg.ChannelID, data)
}

// EmbedReplyNoPing replies to a message with the provided content and/or embed,
// and suppresses all mentions.
func ReplyNoPing(
	st *state.State,
	msg *discord.Message,
	content string,
	embeds ...discord.Embed) (*discord.Message, error) {

	return ReplyNoPingComplex(st, msg, api.SendMessageData{
		Content: content,
		Embeds:  embeds,
	})
}

// EmbedReplyNoPing replies to a message with the provided content and
// suppresses all mentions.
func TextReplyNoPing(
	st *state.State,
	msg *discord.Message,
	content string) (*discord.Message, error) {

	return ReplyNoPing(st, msg, content)
}

// EmbedReplyNoPing replies to a message with the provided embed and suppresses
// all mentions.
func EmbedReplyNoPing(
	st *state.State,
	msg *discord.Message,
	embeds ...discord.Embed) (*discord.Message, error) {

	return ReplyNoPing(st, msg, "", embeds...)
}

// MessageRespond responds to an interaction with the supplied response data.
func MessageRespond(
	st *state.State,
	interaction *discord.InteractionEvent,
	data api.InteractionResponseData) error {

	response := MessageResponse(data)
	return st.RespondInteraction(interaction.ID, interaction.Token, *response)
}

// MessageRespondSimple responds to an interaction with the supplied
// content and embed(s).
func MessageRespondSimple(
	st *state.State,
	interaction *discord.InteractionEvent,
	content string,
	embeds ...discord.Embed) error {

	respData := api.InteractionResponseData{
		Content: option.NewNullableString(content),
		Embeds:  &embeds,
	}

	return MessageRespond(st, interaction, respData)
}

// MessageRespondText responds to an interaction with the supplied
// content.
func MessageRespondText(
	st *state.State,
	interaction *discord.InteractionEvent,
	content string) error {

	return MessageRespondSimple(st, interaction, content)
}

// MessageRespondEmbed responds to an interaction with the
// supplied embed(s).
func MessageRespondEmbed(
	st *state.State,
	interaction *discord.InteractionEvent,
	embeds ...discord.Embed) error {

	return MessageRespondSimple(st, interaction, "", embeds...)
}

// DeferResponse responds to an interaction with a deferred message
// interaction with source response, signalling to the response that
// the full response will come later.
func DeferResponse(
	st *state.State, interaction *discord.InteractionEvent) error {

	return st.RespondInteraction(
		interaction.ID,
		interaction.Token,
		*DeferredMessageResponse(api.InteractionResponseData{}),
	)
}

// FollowupRespond follows up the supplied deferred interaction with the
// supplied response data.
func FollowupRespond(
	st *state.State,
	interaction *discord.InteractionEvent,
	data api.InteractionResponseData) (*discord.Message, error) {

	return st.CreateInteractionFollowup(
		interaction.AppID, interaction.Token, data,
	)
}

// FollowupRespondSimple responds to a deferred interaction with the supplied
// content and embed(s).
func FollowupRespondSimple(
	st *state.State,
	interaction *discord.InteractionEvent,
	content string,
	embeds ...discord.Embed) (*discord.Message, error) {

	respData := api.InteractionResponseData{
		Content: option.NewNullableString(content),
		Embeds:  &embeds,
	}

	return FollowupRespond(st, interaction, respData)
}

// FollowupRespondText responds to a deferred interaction with the supplied
// content.
func FollowupRespondText(
	st *state.State,
	interaction *discord.InteractionEvent,
	content string) (*discord.Message, error) {

	return FollowupRespondSimple(st, interaction, content)
}

// FollowupRespondEmbed responds to a deffered interaction with the
// supplied embed(s).
func FollowupRespondEmbed(
	st *state.State,
	interaction *discord.InteractionEvent,
	embeds ...discord.Embed) (*discord.Message, error) {

	return FollowupRespondSimple(st, interaction, "", embeds...)
}

// ModalRespond responds to an interaction with the supplied response data.
func ModalRespond(
	st *state.State,
	interaction *discord.InteractionEvent,
	data api.InteractionResponseData) error {

	response := ModalResponse(data)
	return st.RespondInteraction(interaction.ID, interaction.Token, *response)
}
