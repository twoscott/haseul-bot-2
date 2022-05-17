// Package bot util provides helper functions pertaining to the bot itself.
package botutil

import (
	"time"

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
