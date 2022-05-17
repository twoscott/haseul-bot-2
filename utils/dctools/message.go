package dctools

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
)

// MessageLink returns a discord URL pointing to the message the provided
// details belong to.
func MessageLink(
	guildID discord.GuildID,
	channelID discord.ChannelID,
	messageID discord.MessageID) string {

	return fmt.Sprintf(
		"https://discord.com/channels/%d/%d/%d",
		guildID, channelID, messageID,
	)
}

// IsUserMessage returns whether the message was sent by a user, and not
// the system.
func IsUserMessage(msgType discord.MessageType) bool {
	return msgType == discord.DefaultMessage ||
		msgType == discord.InlinedReplyMessage ||
		msgType == discord.ApplicationCommandMessage
}
