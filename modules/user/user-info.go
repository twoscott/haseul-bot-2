package user

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/httputil"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var userInfoCommand = &router.SubCommand{
	Name:        "info",
	Description: "Displays information about a Discord user",
	Handler: &router.CommandHandler{
		Executor: userInfoExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.UserOption{
			OptionName:  "user",
			Description: "The user to display information for",
		},
	},
}

func userInfoExec(ctx router.CommandCtx) {
	userSnowflake, _ := ctx.Options.Find("user").SnowflakeValue()

	var member *discord.Member
	userID := discord.UserID(userSnowflake)
	if !userID.IsValid() {
		userID = ctx.Interaction.SenderID()
		member = ctx.Interaction.Member
	}
	if member == nil {
		member, _ = ctx.State.Member(ctx.Interaction.GuildID, userID)
	}

	var user *discord.User
	var err error
	if member == nil {
		user, err = ctx.State.User(userID)
	} else {
		user = &member.User
	}

	if dctools.ErrUnknownUser(err) {
		ctx.RespondWarning("User does not exist.")
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching user data.")
		return
	}

	var embed *discord.Embed
	if member == nil {
		embed = userEmbed(user)
	} else {
		embed = memberEmbed(&ctx, member)
	}

	ctx.RespondEmbed(*embed)
}

func userEmbed(user *discord.User) *discord.Embed {
	embed := discord.Embed{
		Title: user.Tag(),
		Thumbnail: &discord.EmbedThumbnail{
			URL: user.AvatarURL(),
		},
		Description: user.Mention(),
		Fields:      []discord.EmbedField{},
		Color:       dctools.EmbedColour(user.Accent),
		Footer: &discord.EmbedFooter{
			Text: "Member #N/A - User not in server",
		},
	}

	if user.PublicFlags != discord.NoFlag {
		badges := dctools.BadgeEmojiStrings(user.PublicFlags)
		badgeString := strings.Join(badges, util.ThinSpace)

		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  "Badges",
			Value: dctools.MinimiseEmojiString(badgeString),
		})
	}

	embed.Fields = append(embed.Fields, discord.EmbedField{
		Name:   "Account Created",
		Value:  dctools.EmbedTime(user.CreatedAt()),
		Inline: true,
	})

	var avatarUploaded time.Time
	if user.Avatar != "" {
		avatarUploaded, _ = httputil.ImgUploadTime(user.AvatarURL())
	}
	if !avatarUploaded.IsZero() {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:   "Avatar Uploaded",
			Value:  dctools.EmbedTime(avatarUploaded),
			Inline: true,
		})
	}

	embed.Fields = append(embed.Fields, discord.EmbedField{
		Name:  "User ID",
		Value: user.ID.String(),
	})

	if user.Banner != "" {
		url := dctools.ResizeImage(user.BannerURL(), 4096)
		embed.Image = &discord.EmbedImage{URL: url}
	}

	return &embed
}

func memberEmbed(
	ctx *router.CommandCtx, member *discord.Member) *discord.Embed {

	colour, _ := ctx.State.MemberColor(ctx.Interaction.GuildID, member.User.ID)
	if dctools.ColourInvalid(colour) {
		colour = member.User.Accent
	}

	embed := discord.Embed{
		Title: member.User.Tag(),
		Thumbnail: &discord.EmbedThumbnail{
			URL: member.User.AvatarURL(),
		},
		Description: member.Mention(),
		Fields:      []discord.EmbedField{},
		Color:       dctools.EmbedColour(colour),
	}

	var badgeString string
	if member.User.PublicFlags != discord.NoFlag {
		badges := dctools.BadgeEmojiStrings(member.User.PublicFlags)
		badgeString = strings.Join(badges, util.ThinSpace)
	}
	if badgeString != "" {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  "Badges",
			Value: dctools.MinimiseEmojiString(badgeString),
		})
	}

	if member.BoostedSince.IsValid() {
		field := dctools.EmbedTime(member.BoostedSince.Time())

		emoji := dctools.UserBoostEmoji(member.BoostedSince)
		if emoji != nil {
			field = emoji.String() + util.ThinSpace + field
		}

		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  "Boosting Since",
			Value: field,
		})
	}

	embed.Fields = append(embed.Fields, discord.EmbedField{
		Name:   "Account Created",
		Value:  dctools.EmbedTime(member.User.CreatedAt()),
		Inline: true,
	})

	if member.Joined.IsValid() {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:   "Joined Server",
			Value:  dctools.EmbedTime(member.Joined.Time()),
			Inline: true,
		})
	}

	if len(member.RoleIDs) > 0 {
		roles := ""
		for _, roleID := range member.RoleIDs {
			mention := roleID.Mention()
			if len(roles)+len(mention) >= 2048 {
				break
			}
			roles += mention + " "
		}

		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  "Roles",
			Value: strings.TrimSpace(roles),
		})
	}

	var avatarUploaded time.Time
	if member.User.Avatar != "" {
		avatarUploaded, _ = httputil.ImgUploadTime(member.User.AvatarURL())
	}
	if !avatarUploaded.IsZero() {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:   "Avatar Uploaded",
			Value:  dctools.EmbedTime(avatarUploaded),
			Inline: true,
		})
	}

	embed.Fields = append(embed.Fields, discord.EmbedField{
		Name:   "User ID",
		Value:  member.User.ID.String(),
		Inline: true,
	})

	if member.User.Banner != "" {
		url := dctools.ResizeImage(member.User.BannerURL(), 4096)
		embed.Image = &discord.EmbedImage{URL: url}
	}

	memberNumber := dctools.MemberNumber(
		ctx.State, ctx.Interaction.GuildID, member,
	)
	if memberNumber < 1 {
		embed.Footer = &discord.EmbedFooter{
			Text: "Member #--",
		}
	} else {
		embed.Footer = &discord.EmbedFooter{
			Text: fmt.Sprintf("Member #%d", memberNumber),
		}
	}

	return &embed
}
