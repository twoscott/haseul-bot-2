// Package bot util provides helper functions pertaining to the bot itself.
package botutil

import (
	"errors"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/twoscott/haseul-bot-2/config"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

const (
	Website = "https://haseulbot.xyz"
	Discord = "https://discord.gg/w4q5qux"
	Patreon = "https://www.patreon.com/haseulbot"
)

var startTime = time.Now()

// Uptime returns a time period between the time the bot started, and now.
func Uptime() *util.TimePeriod {
	return util.TimeDiff(startTime, time.Now())
}

func Log(st *state.State, messageData api.SendMessageData) (*discord.Message, error) {
	logChannelID := config.GetInstance().Discord.LogChannelID
	if !logChannelID.IsValid() {
		return nil, errors.New("invalid log channel to log to")
	}

	return st.SendMessageComplex(logChannelID, messageData)
}

func LogText(st *state.State, content string) (*discord.Message, error) {
	logChannelID := config.GetInstance().Discord.LogChannelID
	if !logChannelID.IsValid() {
		return nil, errors.New("invalid log channel to log to")
	}

	return st.SendMessage(logChannelID, content)
}
