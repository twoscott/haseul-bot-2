package logs

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/botutil"
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

	msg := welcomeMember(rt.State, member, *guild)
	logMemberJoin(rt.State, member, *guild, msg)
}

func welcomeMember(
	st *state.State,
	member discord.Member,
	guild discord.Guild) *discord.Message {

	welcome, err := db.Guilds.WelcomeConfig(guild.ID)
	if err != nil {
		log.Println(err)
		return nil
	}
	if !welcome.ChannelID.IsValid() {
		return nil
	}

	colour := discord.Color(dctools.EmbedBackColour)
	if welcome.Colour() != dctools.EmbedBackColour {
		colour = welcome.Colour()
	} else {
		url := dctools.MemberAvatarURL(member, guild.ID)
		c, err := dctools.EmbedImageColour(url)
		if err == nil {
			colour = c
		}
	}

	embed := discord.Embed{
		Title:       welcome.Title(),
		Description: welcome.FormattedMessage(member, guild),
		Thumbnail: &discord.EmbedThumbnail{
			URL: member.User.AvatarURL(),
		},
		Color: colour,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("Member #%d", guild.ApproximateMembers),
		},
		Timestamp: member.Joined,
	}

	msg, err := st.SendEmbeds(welcome.ChannelID, embed)
	if err != nil {
		log.Println(err)
		return nil
	}

	return msg
}

func logMemberJoin(
	st *state.State,
	member discord.Member,
	guild discord.Guild,
	welcomeMsg *discord.Message) {

	logChannelID, err := db.Guilds.GetMemberLogsChannel(guild.ID)
	if err != nil {
		log.Println(err)
		return
	}
	if !logChannelID.IsValid() {
		return
	}

	var (
		usedInvite  *discord.Invite
		inviteField = "Currently Unavailable"
	)

	canManageGuild, err := botutil.HasPermissions(
		st,
		logChannelID,
		discord.PermissionManageGuild,
	)
	if err != nil || canManageGuild {
		usedInvite, err = inviteTracker.ResolveInvite(st, guild.ID)
		if err != nil {
			log.Println(err)
		}
	} else {
		inviteField = "I require `MANAGE_GUILD` permissions"
	}

	if usedInvite != nil {
		uses := int64(usedInvite.Uses)
		inviteField = fmt.Sprintf(
			"%s (%s)",
			usedInvite.URL(),
			util.PluraliseWithCount("use", uses),
		)

		if usedInvite.Inviter != nil {
			inviter := usedInvite.Inviter.Tag()
			inviteField += fmt.Sprintf("\nCreated by %s", inviter)
		}

		if usedInvite.CreatedAt.IsValid() {
			inviteField += " " + dctools.UnixTimestampStyled(
				usedInvite.CreatedAt.Time(),
				dctools.RelativeTime,
			)
		}
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
				Value:  dctools.UnixTimestamp(member.Joined.Time()),
				Inline: true,
			},
			{
				Name:   "Account Created On",
				Value:  dctools.UnixTimestamp(member.User.CreatedAt()),
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
			Text: fmt.Sprintf(
				"Member #%s",
				humanize.Comma(int64(guild.ApproximateMembers)),
			),
		},
	}

	if discord.HasFlag(
		uint64(member.User.Flags), uint64(discord.LikelySpammer)) {

		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:   "Flags",
			Value:  fmt.Sprintln(util.WarningSymbol, "Likely Spammer"),
			Inline: true,
		})
	}

	var components discord.ContainerComponents
	if welcomeMsg != nil {
		components = discord.Components(
			&discord.ActionRowComponent{
				&discord.ButtonComponent{
					Label: "Jump to Welcome Message",
					Style: discord.LinkButtonStyle(welcomeMsg.URL()),
				},
			},
		)
	}

	st.SendMessageComplex(logChannelID, api.SendMessageData{
		Embeds:     []discord.Embed{embed},
		Components: components,
	})
}
