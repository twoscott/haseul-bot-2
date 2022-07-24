package logs

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

func logNewMember(
	rt *router.Router, member discord.Member, guildID discord.GuildID) {

	guild, err := rt.State.Session.GuildWithCount(guildID)
	if err != nil {
		guild, err = rt.State.Guild(guildID)
	}
	if err != nil {
		*guild = discord.Guild{Name: "the server"}
	}

	logMemberJoin(rt, member, *guild)
	welcomeMember(rt, member, *guild)
}

func logMemberJoin(
	rt *router.Router, member discord.Member, guild discord.Guild) {

	logChannelID, err := db.Guilds.MemberLogsChannel(guild.ID)
	if err != nil {
		log.Println(err)
		return
	}
	if !logChannelID.IsValid() {
		return
	}

	usedInvite, err := inviteTracker.ResolveInvite(rt.State, guild.ID)
	if err != nil {
		log.Println(err)
	}

	inviteField := "Unknown"
	if usedInvite != nil {
		inviteField = fmt.Sprintf(
			"%s (%d uses)\nCreated by %s",
			usedInvite.URL(),
			usedInvite.Uses,
			usedInvite.Inviter.Tag(),
		)
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
			{
				Name:   "Invite Used",
				Value:  inviteField,
				Inline: false,
			},
			{
				Name:   "User ID",
				Value:  member.User.ID.String(),
				Inline: true,
			},
		},
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("Member #%d", guild.ApproximateMembers),
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

func welcomeMember(
	rt *router.Router, member discord.Member, guild discord.Guild) {

	welcome, err := db.Guilds.WelcomeConfig(guild.ID)
	if err != nil {
		log.Println(err)
		return
	}
	if !welcome.ChannelID.IsValid() {
		return
	}

	embed := discord.Embed{
		Title:       welcome.Title(),
		Description: welcome.Message(member, guild),
		Thumbnail: &discord.EmbedThumbnail{
			URL: member.User.AvatarURL(),
		},
		Color: welcome.Colour(),
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("Member #%d", guild.ApproximateMembers),
		},
		Timestamp: member.Joined,
	}

	if discord.HasFlag(
		uint64(member.User.Flags), uint64(discord.LikelySpammer)) {

		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:   "Flags",
			Value:  fmt.Sprint(util.WarningSymbol, "Likely Spammer"),
			Inline: true,
		})
	}

	rt.State.SendEmbeds(welcome.ChannelID, embed)
}
