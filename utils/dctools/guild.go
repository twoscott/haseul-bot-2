package dctools

import (
	"slices"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/enescakir/emoji"
)

// MemberNumber returns the position a member holds in the sequence of all guild
// members ordered by their join times.
//
// TODO: should return an error instead of zeroes.
func MemberNumber(
	st *state.State, guildID discord.GuildID, member *discord.Member) int {

	if st == nil || member == nil {
		return 0
	}

	members, _ := st.Session.Members(guildID, 0)
	if len(members) < 1 {
		return 0
	}

	memberNo := 1
	for _, m := range members {
		if m.Joined.Time().Before(member.Joined.Time()) {
			memberNo++
		}
	}

	return memberNo
}

// IsOwner returns whether the user ID is the owner of the guild.
func IsOwner(guild discord.Guild, userID discord.UserID) bool {
	return guild.OwnerID == userID
}

// IsEveryoneRole returns whether the role is the @everyone role for the server.
func IsEveryoneRole(guildID discord.GuildID, roleID discord.RoleID) bool {
	return discord.Snowflake(guildID) == discord.Snowflake(roleID)
}

// GuildStatuses returns the approximate number of online and offline
// members in a guild.
func GuildStatuses(guild discord.Guild) (uint64, uint64) {
	online := guild.ApproximatePresences
	offline := guild.ApproximateMembers - online
	return online, offline
}

// GuildHasFeature returns whether a feature exists in a guild.
func GuildHasFeature(guild discord.Guild, feature discord.GuildFeature) bool {
	return slices.Contains(guild.Features, feature)
}

// GuildEmojiLimit returns the emoji limit for a given guild boost level.
func GuildEmojiLimit(level discord.NitroBoost) int {
	switch level {
	case discord.NitroLevel1:
		return 100
	case discord.NitroLevel2:
		return 150
	case discord.NitroLevel3:
		return 250
	default:
		return 50
	}
}

// GuildEmojiLimit returns the file size upload limit in MB for a given
// guild boost level.
func GuildUploadLimitMB(level discord.NitroBoost) int {
	switch level {
	case discord.NitroLevel2:
		return 50
	case discord.NitroLevel3:
		return 100
	default:
		return 8
	}
}

// GuildAudioQualityKbps returns the audio quality in Kbps for a given
// guild boost level.
func GuildAudioQualityKbps(level discord.NitroBoost) int {
	switch level {
	case discord.NitroLevel1:
		return 128
	case discord.NitroLevel2:
		return 256
	case discord.NitroLevel3:
		return 384
	default:
		return 96
	}
}

// GuildStreamQuality returns a string representing the video stream quality
// for a given guild boost level.
func GuildStreamQuality(level discord.NitroBoost) string {
	switch level {
	case discord.NitroLevel1:
		return "720p 60"
	case discord.NitroLevel2:
		return "1080p 60"
	default:
		return "720p 30"
	}
}

// GuildRegionText returns a nicely formatted representation of a region
// translated from the voice region strings returned from Discord's API.
func GuildRegionText(region string) string {
	switch region {
	case "us-west":
		return emoji.Sprintln(emoji.FlagForUnitedStates, "US West")
	case "vip-us-west":
		return emoji.Sprintln(emoji.FlagForUnitedStates, "VIP US West")
	case "us-east":
		return emoji.Sprintln(emoji.FlagForUnitedStates, "US East")
	case "vip-us-east":
		return emoji.Sprintln(emoji.FlagForUnitedStates, "VIP US East")
	case "us-south":
		return emoji.Sprintln(emoji.FlagForUnitedStates, "US South")
	case "us-central":
		return emoji.Sprintln(emoji.FlagForUnitedStates, "US Central")

	case "eu-west":
		return emoji.Sprintln(emoji.FlagForEuropeanUnion, "EU West")
	case "eu-central":
		return emoji.Sprintln(emoji.FlagForEuropeanUnion, "EU Central")
	case "europe":
		return emoji.Sprintln(emoji.FlagForEuropeanUnion, "Europe")

	case "amsterdam":
		return emoji.Sprintln(emoji.FlagForNetherlands, "Amsterdam")
	case "vip-amsterdam":
		return emoji.Sprintln(emoji.FlagForNetherlands, "VIP Amsterdam")

	case "singapore":
		return emoji.Sprintln(emoji.FlagForSingapore, "Singapore")
	case "london":
		return emoji.Sprintln(emoji.FlagForUnitedKingdom, "London")
	case "sydney":
		return emoji.Sprintln(emoji.FlagForAustralia, "Sydney")
	case "frankfurt":
		return emoji.Sprintln(emoji.FlagForGermany, "Frankfurt")
	case "brazil":
		return emoji.Sprintln(emoji.FlagForBrazil, "Brazil")
	case "hongkong":
		return emoji.Sprintln(emoji.FlagForHongKongSarChina, "Hong Kong")
	case "russia":
		return emoji.Sprintln(emoji.FlagForRussia, "Russia")
	case "japan":
		return emoji.Sprintln(emoji.FlagForJapan, "Japan")
	case "southafrica":
		return emoji.Sprintln(emoji.FlagForSouthAfrica, "South Africa")
	case "south-korea":
		return emoji.Sprintln(emoji.FlagForSouthKorea, "South Korea")
	case "india":
		return emoji.Sprintln(emoji.FlagForIndia, "India")
	case "dubai":
		return emoji.Sprintln(emoji.FlagForUnitedArabEmirates, "Dubai")
	default:
		return emoji.Sprintln(emoji.GlobeShowingEuropeAfrica, emoji.QuestionMark)
	}
}

// FindRoleIDByName returns a role ID of a role in roles where its name is
// roleName.
func FindRoleIDByName(roles []discord.Role, roleName string) discord.RoleID {
	for _, role := range roles {
		if role.Name == roleName {
			return role.ID
		}
	}

	return discord.NullRoleID
}
