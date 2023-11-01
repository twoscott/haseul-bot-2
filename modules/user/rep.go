package user

import (
	"log"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/database/repdb"
	"github.com/twoscott/haseul-bot-2/router"
)

const streakEndingHours = 3 * time.Hour

var repCommand = &router.Command{
	Name:        "rep",
	Description: "Commands pertaining to user reps",
}

func getStreakEmojiString(streak repdb.RepStreak) string {
	emojis := []string{}
	days := streak.Days()

	if days >= 1 {
		emojis = append(emojis, "🔥")
	}
	if days >= 100 {
		emojis = append(emojis, "💯")
	}

	user1BF := getStreakBestFriend(streak.UserID1)
	user2BF := getStreakBestFriend(streak.UserID2)

	if user1BF == streak.UserID2 && user2BF == streak.UserID1 {
		bfEmoji := getBestFriendEmoji(days)
		if bfEmoji != "" {
			emojis = append(emojis, bfEmoji)
		}
	}

	topEmoji := getTopStreakEmoji(streak)
	if topEmoji != "" {
		emojis = append(emojis, topEmoji)
	}

	expTime, err := db.Reps.GetTimeToStreakExpiry(streak)
	if err != nil {
		log.Println(err)
	} else if expTime < streakEndingHours {
		emojis = append(emojis, "⌛")
	}

	return strings.Join(emojis, " ")
}

func getStreakBestFriend(userID discord.UserID) discord.UserID {
	streaks, err := db.Reps.GetUserStreaks(userID)
	if err != nil || len(streaks) == 0 {
		log.Println(err)
		return discord.NullUserID
	}

	bfStreak := streaks[0]
	for _, s := range streaks[1:] {
		if s.FirstRep.Unix() < bfStreak.FirstRep.Unix() {
			bfStreak = s
		}
	}

	return bfStreak.OtherUser(userID)
}

func getBestFriendEmoji(days int) string {
	switch {
	case days >= 1000:
		return "💜"
	case days >= 365:
		return "💞"
	case days >= 90:
		return "💕"
	case days >= 30:
		return "💖"
	case days >= 14:
		return "❤️"
	case days >= 7:
		return "💛"
	default:
		return ""
	}
}

func getTopStreakEmoji(streak repdb.RepStreak) string {
	topStreaks, err := db.Reps.GetTopStreaks(3)
	if err != nil {
		log.Println(err)
	}

	topRankEmojis := []string{"🌟", "⭐", "✨"}
	for i, ts := range topStreaks {
		if streak.Equals(ts) {
			return topRankEmojis[i]
		}
	}

	return ""
}
