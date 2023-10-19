package dctools

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var emojiRegex = regexp.MustCompile(`<(a?):(\S+?):(\d+)>`)

// ParseEmoji parses and returns an emoji from an emoji string
func ParseEmoji(emojiString string) (*discord.Emoji, error) {
	if emojiString == "" {
		return nil, errors.New("no emoji string provided")
	}

	match := emojiRegex.FindStringSubmatch(emojiString)
	if match == nil {
		return nil, errors.New("invalid emoji provided")
	}

	emojiID, err := strconv.ParseUint(match[3], 10, 64)
	if err != nil {
		return nil, err
	}

	emoji := discord.Emoji{
		ID:   discord.EmojiID(emojiID),
		Name: match[2],
	}

	if match[1] == "a" {
		emoji.Animated = true
	}

	return &emoji, nil
}

// BadgeEmojiStrings returns a string of emojis pertaining to the passed
// user flags.
func BadgeEmojiStrings(flags discord.UserFlags) []string {
	if flags == discord.NoFlag {
		return nil
	}

	emojiMap := map[discord.UserFlags]discord.Emoji{
		discord.Employee: {
			ID: 844332396690145320, Name: "Employee"},
		discord.Partner: {
			ID: 844332396782026763, Name: "Partner"},
		discord.HypeSquadEvents: {
			ID: 844332396890292295, Name: "HypeSquadEvents"},
		discord.BugHunterLvl1: {
			ID: 844332396488949800, Name: "BugHunterLvl1"},
		discord.HouseBravery: {
			ID: 844332396983484426, Name: "HouseBravery"},
		discord.HouseBrilliance: {
			ID: 844332396593545218, Name: "HouseBrilliance"},
		discord.HouseBalance: {
			ID: 844332396865650698, Name: "HouseBalance"},
		discord.EarlySupporter: {
			ID: 844332396584894474, Name: "EarlySupporter"},
		discord.BugHunterLvl2: {
			ID: 844332396413583401, Name: "BugHunterLvl2"},
		discord.VerifiedBotDeveloper: {
			ID: 844332396735627326, Name: "VerifiedBotDeveloper"},

		// TODO: add badges https://discord.com/developers/docs/resources/user#user-object-user-flags
	}

	emojis := make([]string, 0, len(emojiMap))
	for f := discord.Employee; f <= discord.VerifiedBotDeveloper; f <<= 1 {
		if discord.HasFlag(uint64(flags), uint64(f)) {
			e := emojiMap[f]
			emojis = append(emojis, e.String())
		}
	}

	return emojis
}

// UserBoostEmoji returns an emoji of the user boost badge reflecing how long
// the user has been boosting the server.
func UserBoostEmoji(since discord.Timestamp) *discord.Emoji {
	if !since.IsValid() {
		return nil
	}

	period := util.TimeDiff(since.Time(), time.Now())
	months := (period.Years * 12) + period.Months
	if months < 1 {
		return nil
	}

	switch {
	case months < 2:
		return &discord.Emoji{ID: 844332422563627008, Name: "BoostingLvl1"}
	case months < 3:
		return &discord.Emoji{ID: 844332422526402581, Name: "BoostingLvl2"}
	case months < 6:
		return &discord.Emoji{ID: 844332422484983808, Name: "BoostingLvl3"}
	case months < 9:
		return &discord.Emoji{ID: 844332422610157568, Name: "BoostingLvl4"}
	case months < 12:
		return &discord.Emoji{ID: 844332422556549150, Name: "BoostingLvl5"}
	case months < 15:
		return &discord.Emoji{ID: 844332422589448212, Name: "BoostingLvl6"}
	case months < 18:
		return &discord.Emoji{ID: 844332422485114901, Name: "BoostingLvl7"}
	case months < 24:
		return &discord.Emoji{ID: 844332422706233354, Name: "BoostingLvl8"}
	default:
		return &discord.Emoji{ID: 844332422689456168, Name: "BoostingLvl9"}
	}
}

// BoostLevelEmoji returns the appropriate server boost level emoji based on the
// level provided.
func BoostLevelEmoji(level discord.NitroBoost) *discord.Emoji {
	switch level {
	case discord.NitroLevel1:
		return &discord.Emoji{ID: 846259472183328768, Name: "ServerLevel1"}
	case discord.NitroLevel2:
		return &discord.Emoji{ID: 846259471868887050, Name: "ServerLevel2"}
	case discord.NitroLevel3:
		return &discord.Emoji{ID: 846259471931932682, Name: "ServerLevel3"}
	default:
		return &discord.Emoji{ID: 846259471646326824, Name: "ServerLevel0"}
	}
}

// PartneredEmoji returns an emoji of the partnered server symbol.
func PartneredEmoji() *discord.Emoji {
	return &discord.Emoji{ID: 845863558751322173, Name: "Partnered"}
}

// OnlineEmoji returns an emoji of the online presence.
func OnlineEmoji() *discord.Emoji {
	return &discord.Emoji{ID: 846209270647488512, Name: "StatusOnline"}
}

// OfflineEmoji returns an emoji of the offline presence.
func OfflineEmoji() *discord.Emoji {
	return &discord.Emoji{ID: 846209270230810636, Name: "StatusOffline"}
}

// BoostersEmoji returns the emoji for server boosters.
func BoostersEmoji() *discord.Emoji {
	return &discord.Emoji{ID: 846260731099086858, Name: "Boosters"}
}

// ButtonCheckEmoji returns a plain white check emoji for use with buttons.
func ButtonCheckEmoji() *discord.Emoji {
	return &discord.Emoji{ID: 847992521757949992, Name: "ButtonCheck"}
}

// MinimiseEmojiString appends an invisible character that forces the emojis to
// be parsed as if they were a part of a string of text characters.
func MinimiseEmojiString(emojis string) string {
	return emojis + "\u2800"
}
