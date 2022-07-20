package logs

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

func logMemberJoin(
	rt *router.Router, member discord.Member, guildID discord.GuildID) {

	logChannelID, err := db.Guilds.GetMemberLogsChannelID(guildID)
	if err != nil {
		log.Println(err)
		return
	}
	if !logChannelID.IsValid() {
		return
	}

	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: "User Joined",
		},
		Title:       member.User.Tag(),
		Description: fmt.Sprintln(member.Mention(), "joined the server"),
		Thumbnail: &discord.EmbedThumbnail{
			URL: member.User.AvatarURL(),
		},
		Color: dctools.GreenColour,
		Fields: []discord.EmbedField{
			{
				Name:   "User Joined On",
				Value:  dctools.EmbedTime(member.Joined.Time()),
				Inline: true,
			},
			{
				Name:   "Account Created On",
				Value:  dctools.EmbedTime(member.User.CreatedAt()),
				Inline: true,
			},
			dctools.EmptyEmbedField(),
			{
				Name:   "User ID",
				Value:  member.User.ID.String(),
				Inline: true,
			},
		},
	}

	if discord.HasFlag(
		uint64(member.User.Flags), uint64(discord.LikelySpammer)) {

		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:   "Flags",
			Value:  fmt.Sprint(util.WarningSymbol, "Likely Spammer"),
			Inline: true,
		})
	}

	rt.State.SendEmbeds(logChannelID, embed)
}
