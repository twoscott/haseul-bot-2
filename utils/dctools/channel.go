package dctools

import "github.com/diamondburned/arikawa/v3/discord"

// IsTextChannel returns whether the channel is a text or news channel.
func IsTextChannel(channelType discord.ChannelType) bool {
	return channelType == discord.GuildText || channelType == discord.GuildAnnouncement
}

// GetChannelString returns the channel name prefixed with a # if the channel is
// a text channel.
func GetChannelString(channel discord.Channel) string {
	if IsTextChannel(channel.Type) {
		return "#" + channel.Name
	}

	return channel.Name
}

// TextChannelTypes returns a slice of channel types that can be considered
// "text" channels.
func TextChannelTypes() []discord.ChannelType {
	return []discord.ChannelType{
		discord.GuildText,
		discord.GuildAnnouncement,
	}
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
	case discord.GuildAnnouncement:
		return "News"
	case discord.GuildStore:
		return "Store"
	case discord.GuildAnnouncementThread:
		return "News Thread"
	case discord.GuildPublicThread:
		return "Thread"
	case discord.GuildPrivateThread:
		return "Private Thread"
	case discord.GuildStageVoice:
		return "Voice Stage"
	default:
		return "Unknown"
	}
}
