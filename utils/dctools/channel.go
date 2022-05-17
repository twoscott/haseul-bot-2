package dctools

import "github.com/diamondburned/arikawa/v3/discord"

// IsTextChannel returns whether the channel is a text or news channel.
func IsTextChannel(chType discord.ChannelType) bool {
	return chType == discord.GuildText ||
		chType == discord.GuildNews
}

// ChannelTypeString returns the channel type in string form.
func ChannelTypeString(chType discord.ChannelType) string {
	switch chType {
	case discord.GuildText:
		return "Text"
	case discord.DirectMessage:
		return "DM"
	case discord.GuildVoice:
		return "Voice"
	case discord.GroupDM:
		return "Group DM"
	case discord.GuildCategory:
		return "Category"
	case discord.GuildNews:
		return "News"
	case discord.GuildStore:
		return "Store"
	default:
		return "Unknown"
	}
}
