package dctools

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

// MemberNumber returns the position a member holds in the sequence of all guild
// members ordered by their join times.
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

// GuildHasFeature returns whether a feature exists in an array of
// guild features.
func GuildHasFeature(
	features []discord.GuildFeature, feature discord.GuildFeature) bool {

	for _, f := range features {
		if f == feature {
			return true
		}
	}

	return false
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
		return "🇺🇸 US West"
	case "vip-us-west":
		return "🇺🇸 VIP US West"
	case "us-east":
		return "🇺🇸 US East"
	case "vip-us-east":
		return "🇺🇸 VIP US East"
	case "us-south":
		return "🇺🇸 US South"
	case "us-central":
		return "🇺🇸 US Central"

	case "eu-west":
		return "🇪🇺 EU West"
	case "eu-central":
		return "🇪🇺 EU Central"
	case "europe":
		return "🇪🇺 Europe"

	case "amsterdam":
		return "🇳🇱 Amsterdam"
	case "vip-amsterdam":
		return "🇳🇱 VIP Amsterdam"

	case "singapore":
		return "🇸🇬 Singapore"
	case "london":
		return "🇬🇧 London"
	case "sydney":
		return "🇦🇺 Sydney"
	case "frankfurt":
		return "🇩🇪 Frankfurt"
	case "brazil":
		return "🇧🇷 Brazil"
	case "hongkong":
		return "🇭🇰 Hong Kong"
	case "russia":
		return "🇷🇺 Russia"
	case "japan":
		return "🇯🇵 Japan"
	case "southafrica":
		return "🇿🇦 South Africa"
	case "south-korea":
		return "🇰🇷 South Korea"
	case "india":
		return "🇮🇳 India"
	case "dubai":
		return "🇦🇪 Dubai"
	default:
		return "🌍❓"
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
