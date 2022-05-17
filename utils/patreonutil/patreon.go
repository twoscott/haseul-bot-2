package patreonutil

import (
	"strconv"
	"sync"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/config"
	"github.com/twoscott/patreon-go"
	"golang.org/x/oauth2"
)

// PatreonHelper wraps the patreon client and adds helper functions for the bot.
type PatreonHelper struct {
	client     *patreon.Client
	campaignID string
}

const (
	ActivePatron   = "active_patron"
	FormerPatron   = "former_patron"
	DeclinedPatron = "declined_patron"
)

var (
	p    *PatreonHelper
	once sync.Once
)

// GetPatreonHelper returns the instance of PatreonHelper.
func GetPatreonHelper() *PatreonHelper {
	once.Do(func() {
		cfg := config.GetInstance()
		accessToken := cfg.Patreon.AccessToken
		campaignID := cfg.Patreon.CampaignID

		token := &oauth2.Token{AccessToken: accessToken}
		tokenSource := oauth2.StaticTokenSource(token)
		httpClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
		patreonClient := patreon.NewClient(httpClient)

		p = &PatreonHelper{
			client:     patreonClient,
			campaignID: campaignID,
		}
	})

	return p
}

// GetPatrons returns all patrons for the campaign ID.
func (p *PatreonHelper) GetPatrons() ([]*patreon.Member, error) {
	members, err := p.client.GetAllMembersByCampaignID(
		p.campaignID,
		patreon.WithIncludes("user", "currently_entitled_tiers"),
		patreon.WithFields(
			"member",
			"patron_status",
			"pledge_relationship_start",
			"currently_entitled_amount_cents",
		),
		patreon.WithFields("user", "social_connections"),
		patreon.WithPageSize(20),
	)
	if err != nil {
		return nil, err
	}

	return members, nil
}

// GetActivePatrons returns all patrons for the campaign ID that are active.
func (p *PatreonHelper) GetActivePatrons() ([]*patreon.Member, error) {
	patrons, err := p.GetPatrons()
	if err != nil {
		return nil, err
	}

	activePatrons := make([]*patreon.Member, 0, len(patrons))
	for _, patron := range patrons {
		if patron.PatronStatus == ActivePatron {
			activePatrons = append(activePatrons, patron)
		}
	}

	return activePatrons, nil
}

// GetActiveDiscordPatron returns a patreon member whose account is linked to
// the discord account with the given user ID.
func (p *PatreonHelper) GetActiveDiscordPatron(
	userID discord.UserID) (*patreon.Member, error) {

	patrons, err := p.GetActivePatrons()
	if err != nil {
		return nil, err
	}

	for _, patron := range patrons {
		dConn := patron.User.SocialConnections.Discord
		if dConn == nil {
			continue
		}

		patronUserID, _ := strconv.ParseInt(dConn.UserID, 10, 64)
		if userID == discord.UserID(patronUserID) {
			return patron, nil
		}
	}

	return nil, nil
}
