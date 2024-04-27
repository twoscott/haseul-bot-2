package inviteutil

import (
	"log"
	"sync"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/database/invitedb"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

// Tracker is used for tracking and resolving used guild invites.
type Tracker struct {
	*database.DB
}

var (
	tracker *Tracker
	once    sync.Once
)

// GetTracker returns the instance of Tracker.
func GetTracker() *Tracker {
	once.Do(func() {
		tracker = &Tracker{
			database.GetInstance(),
		}
	})

	return tracker
}

func (tr *Tracker) trackNewInvites(
	guildID discord.GuildID, invites []discord.Invite) error {
	if len(invites) < 1 {
		return nil
	}

	return tr.Invites.AddAll(guildID, invites)
}

func (tr *Tracker) vanityInvite(
	st *state.State, guildID discord.GuildID) (*discord.Invite, error) {

	guild, err := st.Guild(guildID)
	if err != nil {
		return nil, err
	}

	if !dctools.GuildHasFeature(guild, discord.VanityURL) {
		return nil, nil
	}

	vanity, err := st.Session.GuildVanityInvite(guild.ID)
	if err != nil {
		return nil, err
	}

	// vanity.Code will be empty if Vanity URL is not set for the guild
	// https://discord.com/developers/docs/resources/guild#get-guild-vanity-url
	if vanity.Code == "" {
		return nil, nil
	}

	return vanity, err
}

// ResolveInvite resolves which invite was used by the newest member to join
// the provided guild. If it cannot be determined, the returned invite will be
// nil.
func (tr *Tracker) ResolveInvite(
	st *state.State, guildID discord.GuildID) (*discord.Invite, error) {

	storedInvs, err := tr.Invites.GetAllByGuild(guildID)
	if err != nil {
		return nil, err
	}

	newInvs, err := st.Session.GuildInvites(guildID)
	if err != nil {
		return nil, err
	}
	newVanity, err := tr.vanityInvite(st, guildID)
	if err == nil && newVanity != nil {
		newInvs = append(newInvs, *newVanity)
	}

	for _, oldInv := range storedInvs {
		newInv := findNewInvite(newInvs, oldInv.Code)
		if newInv == nil {
			tr.Invites.Remove(oldInv.Code)
		}
	}

	usedInvites := make([]discord.Invite, 0)
	for _, newInv := range newInvs {
		oldUses := 0

		oldInv := findOldInvite(storedInvs, newInv.Code)
		if oldInv != nil {
			oldUses = oldInv.Uses
		}

		// new invites with 0 uses remain untracked until they are used.
		if newInv.Uses > oldUses {
			usedInvites = append(usedInvites, newInv)
		}
	}

	err = tr.trackNewInvites(guildID, usedInvites)
	if err != nil {
		log.Println("DB Err:", err)
	}

	// either zero or more than one invite has been used since the tracked
	// invites were last updated, thus the invite used is indeterminable.
	if len(usedInvites) != 1 {
		return nil, nil
	}

	return &usedInvites[0], nil
}

func findOldInvite(
	invites []invitedb.Invite, code string) *invitedb.Invite {

	for _, inv := range invites {
		if inv.Code == code {
			return &inv
		}
	}

	return nil
}

func findNewInvite(
	invites []discord.Invite, code string) *discord.Invite {

	for _, inv := range invites {
		if inv.Code == code {
			return &inv
		}
	}

	return nil
}
