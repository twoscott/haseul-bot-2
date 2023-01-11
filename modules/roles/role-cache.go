package roles

import (
	"log"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
)

type roleSelection struct {
	selectedAt time.Time
	roleIDs    []discord.RoleID
}

func (s roleSelection) Age() time.Duration {
	return time.Since(s.selectedAt)
}

type selectMenuTracker struct {
	guildID   discord.GuildID
	messageID discord.MessageID
	userID    discord.UserID
}

type roleCache struct {
	selections map[selectMenuTracker]roleSelection
	maxAge     time.Duration
}

// SetSelection sets the currently selected
func (c *roleCache) SetSelection(
	interaction *discord.InteractionEvent, roleIDs []discord.RoleID) {

	menu := interactionToMenuTracker(*interaction)

	c.selections[menu] = roleSelection{
		selectedAt: time.Now(),
		roleIDs:    roleIDs,
	}
}

// ClearSelection clears the currently selected roles for an interaction.
func (c *roleCache) ClearSelection(interaction *discord.InteractionEvent) {
	menu := interactionToMenuTracker(*interaction)
	delete(c.selections, menu)
}

// GetSelectedRoleIDs gets the roles a user has selected for a specific role
// picker in a guild.
func (c roleCache) GetSelectedRoleIDs(
	interaction *discord.InteractionEvent) []discord.RoleID {

	menu := interactionToMenuTracker(*interaction)

	selection, ok := c.selections[menu]
	if !ok {
		return nil
	}

	return selection.roleIDs
}

// ClearCache clears any selections older than
func (c *roleCache) ClearCache() {
	deleted := 0
	for menu, selection := range c.selections {
		if selection.Age() > c.maxAge {
			delete(c.selections, menu)
			deleted++
		}
	}

	log.Printf("Deleted %d role selections from the cache\n", deleted)
}

// ClearJob starts a job that clears the cache at the provided interval.
func (c *roleCache) ClearJob(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		c.ClearCache()
	}
}

func newRoleCache(maxAge time.Duration) *roleCache {
	rc := roleCache{
		selections: make(map[selectMenuTracker]roleSelection),
		maxAge:     maxAge,
	}

	return &rc
}

func interactionToMenuTracker(
	interaction discord.InteractionEvent) selectMenuTracker {
	return selectMenuTracker{
		guildID:   interaction.GuildID,
		messageID: interaction.Message.ID,
		userID:    interaction.SenderID(),
	}
}
