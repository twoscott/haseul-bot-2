// Package bot util provides helper functions pertaining to the bot itself.
package botutil

import (
	"errors"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/twoscott/haseul-bot-2/config"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

const (
	Discord = "https://discord.gg/w4q5qux"
	Invite  = "https://discord.com/api/oauth2/authorize?client_id=457640781558054925&permissions=517879491681&scope=bot"
)

var startTime = time.Now()

// Uptime returns a time period between the time the bot started, and now.
func Uptime() util.TimePeriod {
	return util.TimeDiff(startTime, time.Now())
}

// Log logs a message in the bot's logging Discord channel.
func Log(st *state.State, messageData api.SendMessageData) (*discord.Message, error) {
	logChannelID := config.GetInstance().Bot.LogChannelID
	if !logChannelID.IsValid() {
		return nil, errors.New("invalid log channel to log to")
	}

	return st.SendMessageComplex(logChannelID, messageData)
}

// LogText logs text in the bot's kiggubg Discord channel.
func LogText(st *state.State, content string) (*discord.Message, error) {
	data := api.SendMessageData{
		Content: content,
	}

	return Log(st, data)
}

// HasAnyPermissions returns whether the bot has any of the provided
// permission (including admin)
func HasAnyPermissions(
	st *state.State,
	channelID discord.ChannelID,
	requiredPerms discord.Permissions) (bool, error) {

	bot, err := st.Me()
	if err != nil {
		return false, err
	}

	permissions, err := st.Permissions(channelID, bot.ID)
	if err != nil {
		return false, err
	}

	return dctools.HasAnyPermOrAdmin(permissions, requiredPerms), nil
}

// HasPermissions returns whether the bot has the permissions provided
// (including admin)
func HasPermissions(
	st *state.State,
	channelID discord.ChannelID,
	requiredPerms discord.Permissions) (bool, error) {

	bot, err := st.Me()
	if err != nil {
		return false, err
	}

	permissions, err := st.Permissions(channelID, bot.ID)
	if err != nil {
		return false, err
	}

	return permissions.Has(requiredPerms), nil
}
