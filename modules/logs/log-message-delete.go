package logs

import (
	"fmt"
	"log"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

func logMessageDelete(
	rt *router.Router, deletedMsg discord.Message) {

	if deletedMsg.Author.Bot {
		return
	}
	if !dctools.IsUserMessage(deletedMsg.Type) {
		return
	}

	logChannelID, err := db.Guilds.GetMessageLogsChannel(deletedMsg.GuildID)
	if err != nil {
		log.Println(err)
		return
	}
	if !logChannelID.IsValid() {
		return
	}

	channelName := "Unknown"
	channel, err := rt.State.Channel(deletedMsg.ChannelID)
	if err == nil {
		channelName = channel.Name
	}

	var proximityMsg *discord.Message
	recentMsgs, err := rt.State.Messages(deletedMsg.ChannelID, 5)
	if err == nil && len(recentMsgs) > 0 {
		proximityMsg = &recentMsgs[0]
		for _, m := range recentMsgs {
			if m.ID < deletedMsg.ID {
				proximityMsg = &m
				break
			}
		}
	}

	dpURL := dctools.ResizeImage(deletedMsg.Author.AvatarURL(), 64)
	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: deletedMsg.Author.Tag(),
			Icon: dpURL,
		},
		Title: "Message Deleted",
		Color: dctools.RedColour,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("#%s", channelName),
		},
	}

	if deletedMsg.Content != "" {
		if len(deletedMsg.Content) > 1024 {
			deletedMsg.Content = deletedMsg.Content[:1021] + "..."
		}
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  "Content",
			Value: deletedMsg.Content,
		})
	}

	if len(deletedMsg.Attachments) > 0 {
		links := make([]string, 0, len(deletedMsg.Attachments))
		for _, a := range deletedMsg.Attachments {
			if len(deletedMsg.Attachments) > 3 {
				links = append(links, "- "+a.Filename)
				break
			}
			links = append(links, "- "+a.Proxy)
		}
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  "Attachments",
			Value: strings.Join(links, "\n"),
		})
	}

	var components discord.ContainerComponents
	if proximityMsg != nil {
		components = discord.Components(
			&discord.ActionRowComponent{
				&discord.ButtonComponent{
					Label: "Jump to Message Area",
					Style: discord.LinkButtonStyle(proximityMsg.URL()),
				},
			},
		)
	}

	rt.State.SendMessageComplex(logChannelID, api.SendMessageData{
		Embeds:     []discord.Embed{embed},
		Components: components,
	})
}
