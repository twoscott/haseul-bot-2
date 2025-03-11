package logs

import (
	"fmt"
	"log"
	"slices"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

func logMessageUpdate(
	rt *router.Router,
	oldMsg discord.Message,
	newMsg discord.Message,
	_ *discord.Member) {

	if newMsg.Author.Bot {
		return
	}
	if !dctools.IsUserMessage(newMsg.Type) {
		return
	}

	logChannelID, err := db.Guilds.GetMessageLogsChannel(newMsg.GuildID)
	if err != nil {
		log.Println(err)
		return
	}
	if !logChannelID.IsValid() {
		return
	}

	if messageGotEmbedded(oldMsg, newMsg) {
		return
	}

	channelName := "Unknown"
	channel, err := rt.State.Channel(newMsg.ChannelID)
	if err == nil {
		channelName = channel.Name
	}

	dpURL := dctools.ResizeImage(newMsg.Author.AvatarURL(), 64)
	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: newMsg.Author.Tag(),
			Icon: dpURL,
		},
		Title: "Message Edited",
		Color: dctools.YellowColour,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("#%s", channelName),
		},
	}

	if oldMsg.Content != newMsg.Content {
		if len(oldMsg.Content) > 1024 {
			oldMsg.Content = oldMsg.Content[:1021] + "..."
		}
		if len(newMsg.Content) > 1024 {
			newMsg.Content = newMsg.Content[:1021] + "..."
		}

		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  "Old Message",
			Value: oldMsg.Content,
		})
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  "New Message",
			Value: newMsg.Content,
		})
	}

	if len(newMsg.Attachments) < len(oldMsg.Attachments) {
		for _, a := range oldMsg.Attachments {
			if !slices.Contains(newMsg.Attachments, a) {
				embed.Fields = append(embed.Fields, discord.EmbedField{
					Name:  "Deleted Attachment",
					Value: a.Proxy,
				})
				break
			}
		}
	}

	components := discord.Components(
		&discord.ActionRowComponent{
			&discord.ButtonComponent{
				Label: "Jump to Message",
				Style: discord.LinkButtonStyle(newMsg.URL()),
			},
		},
	)

	rt.State.SendMessageComplex(logChannelID, api.SendMessageData{
		Embeds:     []discord.Embed{embed},
		Components: components,
	})
}

func messageGotEmbedded(oldMsg, newMsg discord.Message) bool {
	return oldMsg.Content == newMsg.Content && len(oldMsg.Attachments) == len(newMsg.Attachments)
}
