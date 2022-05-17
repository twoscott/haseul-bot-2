// Package dctools provides helper functions and constants relevant to Discord
// and its API.
package dctools

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

// NoMentions is a zero value AllowedMentions object which allows no mentions.
var NoMentions = new(api.AllowedMentions)

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

// ReplyWithSuccess replies to a message with the provided content, prepended
// with a check emoji.
func ReplyWithSuccess(
	st *state.State,
	msg *discord.Message,
	content string) (*discord.Message, error) {

	return TextReplyNoPing(st, msg, Success(content))
}

// SendSuccess sends a message with the provided content, prepended
// with a check emoji.
func SendSuccess(
	st *state.State,
	channelID discord.ChannelID,
	content string) (*discord.Message, error) {

	return st.SendMessage(channelID, Success(content))
}

// ReplyWithError replies to a message with the provided content, prepended
// with a cross emoji.
func ReplyWithError(
	st *state.State,
	msg *discord.Message,
	content string) (*discord.Message, error) {

	return TextReplyNoPing(st, msg, Error(content))
}

// SendError sends a message with the provided content, prepended
// with a cross emoji.
func SendError(
	st *state.State,
	channelID discord.ChannelID,
	content string) (*discord.Message, error) {

	return st.SendMessage(channelID, Error(content))
}

// ReplyWithWarning replies to a message with the provided content, prepended
// with a warning emoji.
func ReplyWithWarning(
	st *state.State,
	msg *discord.Message,
	content string) (*discord.Message, error) {

	return TextReplyNoPing(st, msg, Warning(content))
}

// SendWarning sends a message with the provided content, prepended
// with a warning emoji.
func SendWarning(
	st *state.State,
	channelID discord.ChannelID,
	content string) (*discord.Message, error) {

	return st.SendMessage(channelID, Warning(content))
}
