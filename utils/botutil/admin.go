package botutil

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"golang.org/x/exp/slices"
)

// BotAuthorID returns the user ID of the bot author.
func AuthorID() discord.UserID {
	return 125414437229297664
}

// IsBotAdmin returns whether the given user ID matches the ID of a bot admin.
func IsBotAdmin(userID discord.UserID) bool {
	adminIDs := []discord.UserID{125414437229297664}

	return slices.Contains(adminIDs, userID)
}
