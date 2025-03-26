package botutil

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/config"
)

// BotAuthorID returns the user ID of the bot author.
func AuthorID() discord.UserID {
	cfg := config.GetInstance()
	return cfg.Bot.AdminUserID
}

// IsBotAdmin returns whether the given user ID matches the ID of a bot admin.
func IsBotAdmin(userID discord.UserID) bool {
	cfg := config.GetInstance()
	return userID == cfg.Bot.AdminUserID
}
