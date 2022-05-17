package dctools

import (
	"regexp"
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
)

var (
	channelIDRegex = regexp.MustCompile(`^(?:<#)?(\d+)>?$`)
	emojiIDRegex   = regexp.MustCompile(`<(?:a?:\S+?:)?(\d+)>`)
	roleIDRegex    = regexp.MustCompile(`^(?:<@&)?(\d+)>?$`)
	userIDRegex    = regexp.MustCompile(`^(?:<@!?)?(\d+)>?$`)
)

// ParseChannelID parses and returns a channel ID out of either a channel
// mention or a raw string channel ID.
func ParseChannelID(channel string) discord.ChannelID {
	if channel == "" {
		return discord.NullChannelID
	}

	match := channelIDRegex.FindStringSubmatch(channel)
	if match == nil {
		return discord.NullChannelID
	}

	channelID, err := strconv.ParseUint(match[1], 10, 64)
	if err != nil {
		return discord.NullChannelID
	}

	return discord.ChannelID(channelID)
}

// ParseEmojiID parses and returns an emoji ID out of either an emoji string
// or a raw emoji ID.
func ParseEmojiID(emoji string) discord.EmojiID {
	if emoji == "" {
		return discord.NullEmojiID
	}

	match := emojiIDRegex.FindStringSubmatch(emoji)
	if match == nil {
		return discord.NullEmojiID
	}

	emojiID, err := strconv.ParseUint(match[3], 10, 64)
	if err != nil {
		return discord.NullEmojiID
	}

	return discord.EmojiID(emojiID)
}

// ParseGuildID parses and returns a guild ID from string guild ID.
func ParseGuildID(guild string) discord.GuildID {
	if guild == "" {
		return discord.NullGuildID
	}

	guildID, err := strconv.ParseUint(guild, 10, 64)
	if err != nil {
		return discord.NullGuildID
	}

	return discord.GuildID(guildID)
}

// ParseRoleID parses and returns a role ID out of either a role mention or
// a raw string role ID.
func ParseRoleID(role string) discord.RoleID {
	if role == "" {
		return discord.NullRoleID
	}

	match := roleIDRegex.FindStringSubmatch(role)
	if match == nil {
		return discord.NullRoleID
	}

	roleID, err := strconv.ParseUint(match[1], 10, 64)
	if err != nil {
		return discord.NullRoleID
	}

	return discord.RoleID(roleID)
}

// ParseUserID parses and returns a user ID out of either a user mention or
// a raw string user ID.
func ParseUserID(user string) discord.UserID {
	if user == "" {
		return discord.NullUserID
	}

	match := userIDRegex.FindStringSubmatch(user)
	if match == nil {
		return discord.NullUserID
	}

	userID, err := strconv.ParseUint(match[1], 10, 64)
	if err != nil {
		return discord.NullUserID
	}

	return discord.UserID(userID)
}
