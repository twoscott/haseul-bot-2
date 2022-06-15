package twitter

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
	"golang.org/x/exp/slices"
)

const (
	twitterIcon   = "https://abs.twimg.com/icons/apple-touch-icon-192x192.png"
	twitterColour = 0x1da1f2
)

var URLRegex = regexp.MustCompile("https://twitter.com/([^/]+)/?.*")

func parseUserIfURL(user string) string {
	if user == "" {
		return user
	}

	match := URLRegex.FindStringSubmatch(user)
	if match == nil {
		if strings.HasPrefix(user, "@") {
			user = user[1:]
		}
		return user
	}

	return match[1]
}

func fetchUser(screenName string) (*twitter.User, router.CmdResponse) {
	user, resp, err := twt.Users.Show(&twitter.UserShowParams{
		ScreenName: screenName,
	})
	switch err.(type) {
	case twitter.APIError:
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, router.Warningf(
				"I could not find a user named @%s.", screenName,
			)
		case http.StatusForbidden:
			return nil, router.Warningf(
				"@%s is either private or suspended.", screenName,
			)
		default:
			log.Println(err)
			return nil, router.Errorf(
				"Unknown error occurred while trying to find @%s.",
				screenName,
			)
		}
	}
	if err != nil {
		return nil, router.Errorf(
			"Unknown error occurred while trying to find @%s.", screenName,
		)
	}

	return user, nil
}

func parseChannelArg(
	ctx router.CommandCtx,
	channelID discord.ChannelID) (*discord.Channel, router.CmdResponse) {

	if !channelID.IsValid() {
		return nil, router.Warningf("Malformed Discord channel provided.")
	}

	channel, err := ctx.State.Channel(channelID)
	if dctools.ErrMissingAccess(err) {
		return nil, router.Warningf("I cannot access this channel.")
	}
	if err != nil {
		return nil, router.Warningf("Invalid Discord channel provided.")
	}
	if channel.GuildID != ctx.Interaction.GuildID {
		return nil, router.Warningf(
			"Channel provided must belong to this server.",
		)
	}
	if !dctools.IsTextChannel(channel.Type) {
		return nil, router.Warningf("Channel provided must be a text channel.")
	}

	return channel, nil
}

func dbTwitterCompleter(ctx router.AutocompleteCtx) {
	username := ctx.Options.Find("twitter").String()

	dbUsers, err := db.Twitter.GetUsersByGuild(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		return
	}

	userIDs := make([]int64, 0, len(dbUsers))
	for _, u := range dbUsers {
		userIDs = append(userIDs, u.ID)
	}

	users, resp, err := twt.Users.Lookup(&twitter.UserLookupParams{
		UserID: userIDs,
	})
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println(err)
		return
	}

	usernames := make([]string, 0, len(users))
	for _, u := range users {
		usernames = append(usernames, u.ScreenName)
	}

	var choices api.AutocompleteStringChoices
	if username == "" {
		usernames := slices.Compact(usernames)
		choices = dctools.MakeStringChoices(usernames)
	} else {
		matches := util.SearchSort(usernames, username)
		choices = dctools.MakeStringChoices(matches)
	}

	ctx.RespondChoices(choices)
}
