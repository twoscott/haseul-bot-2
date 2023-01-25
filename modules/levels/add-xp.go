package levels

import (
	"log"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

const baseXp int64 = 5

func addXP(rt *router.Router, msg discord.Message, member *discord.Member) {
	xp := baseXp

	words := strings.Fields(msg.Content)
	wordCount := len(words)

	if wordCount >= 16 {
		xp += 10
	} else if wordCount >= 4 {
		xp += 5
	}

	if len(msg.Attachments) >= 1 {
		xp += 10
	}

	var recentMessage *discord.Message
	channelMsgs, _ := rt.State.MessageStore.Messages(msg.ChannelID)
	for _, m := range channelMsgs {
		if m.ID != msg.ID && m.Author.ID == msg.Author.ID {
			recentMessage = &m
			break
		}
	}

	if recentMessage != nil {
		sentTime := recentMessage.Timestamp.Time()
		if time.Since(sentTime) < (15 * time.Second) {
			xp /= 5
		}
	}

	_, err := db.Levels.AddUserXP(msg.GuildID, msg.Author.ID, xp)
	if err != nil {
		log.Println(err)
	}
}
