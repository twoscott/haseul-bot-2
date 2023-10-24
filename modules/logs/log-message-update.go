package logs

import (
	"fmt"
	"log"

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

	logChannelID, err := db.Guilds.GetMessageLogsChannel(oldMsg.GuildID)
	if err != nil {
		log.Println(err)
		return
	}
	if !logChannelID.IsValid() {
		return
	}

	channelName := "Unknown"
	channel, err := rt.State.Channel(oldMsg.ChannelID)
	if err == nil {
		channelName = channel.Name
	}

	dpURL := dctools.ResizeImage(oldMsg.Author.AvatarURL(), 64)
	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: oldMsg.Author.Tag(),
			Icon: dpURL,
		},
		Title: "Message Edited",
		Color: dctools.YellowColour,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("#%s", channelName),
		},
	}

	if oldMsg.Content != "" {
		if len(oldMsg.Content) > 1024 {
			oldMsg.Content = oldMsg.Content[:1021] + "..."
		}
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  "Old Message",
			Value: oldMsg.Content,
		})
	}

	if newMsg.Content != "" {
		if len(newMsg.Content) > 1024 {
			newMsg.Content = newMsg.Content[:1021] + "..."
		}
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  "New Message",
			Value: newMsg.Content,
		})
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
