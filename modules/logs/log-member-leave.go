package logs

import (
	"fmt"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

func logMemberLeave(
	rt *router.Router, user discord.User, guildID discord.GuildID) {

	logChannelID, err := db.Guilds.MemberLogsChannel(guildID)
	if err != nil {
		log.Println(err)
		return
	}
	if !logChannelID.IsValid() {
		return
	}

	leftAt := time.Now()

	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: "User Left",
		},
		Title:       user.Tag(),
		Description: fmt.Sprintln(user.Mention(), "left the server"),
		Thumbnail: &discord.EmbedThumbnail{
			URL: user.AvatarURL(),
		},
		Color: dctools.RedColour,
		Fields: []discord.EmbedField{
			{
				Name:   "User Left On",
				Value:  dctools.UnixTimestamp(leftAt),
				Inline: true,
			},
			{
				Name:   "Account Created On",
				Value:  dctools.UnixTimestamp(user.CreatedAt()),
				Inline: true,
			},
			dctools.EmptyEmbedField(),
			{
				Name:   "User ID",
				Value:  user.ID.String(),
				Inline: true,
			},
		},
	}

	rt.State.SendEmbeds(logChannelID, embed)
}
