package notifications

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/database/notifdb"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

type notificationMatch struct {
	userID  discord.UserID
	keyword string
}

func checkKeywords(
	rt *router.Router, msg discord.Message, _ *discord.Member) {

	if len(msg.Content) < 1 {
		return
	}

	notifs, err := db.Notifications.GetAllChecking(
		msg.Author.ID, msg.GuildID, msg.ChannelID,
	)
	if err != nil {
		log.Println(err)
		return
	}

	matchChan := make(chan notificationMatch)
	go checkMatches(matchChan, notifs, msg.Content)
	go sendNotifications(rt, matchChan, msg)
}

func checkMatches(
	matchChan chan<- notificationMatch,
	notifs []notifdb.Notification,
	content string) {

	defer close(matchChan)

	var wg sync.WaitGroup
	for _, noti := range notifs {
		wg.Add(1)
		go func(noti notifdb.Notification) {
			defer wg.Done()
			checkMatch(matchChan, noti, content)
		}(noti)
	}

	wg.Wait()
}

func checkMatch(
	matchChan chan<- notificationMatch,
	noti notifdb.Notification,
	content string) {

	rgxString := regexp.QuoteMeta(noti.Keyword)

	switch noti.Type {
	case notifdb.NormalNotification:
		plural := util.PluralSuffix(noti.Keyword)
		possessive := util.PossessiveSuffix(noti.Keyword)
		rgxString = rgxString + `(?:` + possessive + `|` + plural + `)?`
		rgxString = `(?i)(^|\W)` + rgxString + `($|\W)`
	case notifdb.LenientNotification:
		rgxString = `(?i)` + rgxString
	case notifdb.StrictNotification:
		rgxString = `(^|\s)` + rgxString + `($|\s)`
	}

	rgx, err := regexp.Compile(rgxString)
	if err != nil {
		log.Println(err)
		return
	}

	match := rgx.MatchString(content)
	if !match {
		return
	}

	matchChan <- notificationMatch{
		userID:  noti.UserID,
		keyword: noti.Keyword,
	}
}

func sendNotifications(
	rt *router.Router,
	matchChan <-chan notificationMatch,
	msg discord.Message) {

	userMatchSets := make(map[discord.UserID]map[string]struct{})
	exists := struct{}{}

	for match := range matchChan {
		if _, ok := userMatchSets[match.userID]; !ok {
			userMatchSets[match.userID] = make(map[string]struct{})
		}
		userMatchSets[match.userID][match.keyword] = exists
	}

	for userID, matchSet := range userMatchSets {
		matches := make([]string, len(matchSet))

		i := 0
		for keyword := range matchSet {
			matches[i] = keyword
			i++
		}

		go sendNotification(rt, msg, userID, matches)
	}
}

func sendNotification(
	rt *router.Router,
	msg discord.Message,
	userID discord.UserID,
	matches []string) {

	chString := msg.ChannelID.String()
	channel, err := rt.State.Channel(msg.ChannelID)
	if err == nil {
		chString = dctools.GetChannelString(*channel)
	}

	if !canSeeChannel(rt, *channel, userID) {
		return
	}

	dmChannel, err := rt.State.CreatePrivateChannel(userID)
	if err != nil {
		log.Println(err)
		return
	}

	name := msg.Author.DisplayOrUsername()
	matchString := strings.Join(matches, "`, `")
	content := fmt.Sprintf("💬 %s mentioned `%s`",
		dctools.Bold(name), matchString,
	)

	guild, err := rt.State.Guild(msg.GuildID)
	if err == nil {
		content += fmt.Sprintf(" in %s", dctools.Bold(guild.Name))
	}

	colour, _ := rt.State.MemberColor(msg.GuildID, msg.Author.ID)
	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: msg.Author.Tag(),
			Icon: msg.Author.AvatarURL(),
		},
		Description: msg.Content,
		Footer: &discord.EmbedFooter{
			Text: chString,
		},
		Timestamp: msg.Timestamp,
		Color:     dctools.EmbedColour(colour),
	}

	rt.State.SendMessageComplex(dmChannel.ID, api.SendMessageData{
		Content: content,
		Embeds:  []discord.Embed{embed},
		Components: discord.Components(
			&discord.ActionRowComponent{
				&discord.ButtonComponent{
					Label: "Jump to Message",
					Style: discord.LinkButtonStyle(msg.URL()),
				},
			},
		),
	})
}

func canSeeChannel(
	rt *router.Router, channel discord.Channel, userID discord.UserID) bool {

	switch channel.Type {
	case discord.GuildText,
		discord.GuildVoice,
		discord.GuildStageVoice,
		discord.GuildAnnouncement:

		permissions, err := rt.State.Permissions(channel.ID, userID)
		return err == nil && permissions.Has(discord.PermissionViewChannel)
	case discord.GuildPublicThread,
		discord.GuildAnnouncementThread:

		permissions, err := rt.State.Permissions(channel.ParentID, userID)
		return err == nil && permissions.Has(discord.PermissionViewChannel)
	case discord.GuildPrivateThread:
		_, err := rt.State.ThreadMember(channel.ID, userID)
		return err == nil
	}

	return false
}
