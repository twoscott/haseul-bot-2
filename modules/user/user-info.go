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

	user, err := ctx.State.User(userID)
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
		embed = userEmbed(*user)
	} else {
		embed = memberEmbed(ctx, *member, *user)
	}

	ctx.RespondEmbed(*embed)
}

// TODO: combine into single embed function to minimise repetition
func userEmbed(user discord.User) *discord.Embed {
	url := user.AvatarURLWithType(discord.PNGImage)
	colour, err := dctools.EmbedImageColour(url)
	if err != nil || dctools.ColourInvalid(colour) {
		colour = user.Accent
	}

	embed := discord.Embed{
		Title: fmt.Sprintf("%s (@%s)", user.DisplayOrUsername(), user.Tag()),
		Thumbnail: &discord.EmbedThumbnail{
			URL: user.AvatarURL(),
		},
		Description: user.Mention(),
		Fields:      []discord.EmbedField{},
		Color:       dctools.EmbedColour(colour),
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
		Value:  dctools.UnixTimestamp(user.CreatedAt()),
		Inline: true,
	})

	var avatarUploaded time.Time
	if user.Avatar != "" {
		avatarUploaded, _ = httputil.ImgUploadTime(user.AvatarURL())
	}
	if !avatarUploaded.IsZero() {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:   "Avatar Uploaded",
			Value:  dctools.UnixTimestamp(avatarUploaded),
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
	ctx router.CommandCtx,
	member discord.Member,
	user discord.User) *discord.Embed {

	guildID := ctx.Interaction.GuildID
	colour, ok := ctx.State.MemberColor(guildID, user.ID)
	if !ok || dctools.ColourInvalid(colour) {
		url := dctools.MemberAvatarURL(member, guildID)
		colour, _ = dctools.EmbedImageColour(url)
	}

	embed := discord.Embed{
		Title: fmt.Sprintf("%s (@%s)", user.DisplayOrUsername(), user.Tag()),
		Thumbnail: &discord.EmbedThumbnail{
			URL: dctools.MemberAvatarURL(member, guildID),
		},
		Description: member.Mention(),
		Fields:      []discord.EmbedField{},
		Color:       dctools.EmbedColour(colour),
		Timestamp:   discord.Timestamp(user.CreatedAt()),
	}

	var badgeString string
	if user.PublicFlags != discord.NoFlag {
		badges := dctools.BadgeEmojiStrings(user.PublicFlags)
		badgeString = strings.Join(badges, util.ThinSpace)
	}
	if badgeString != "" {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  "Badges",
			Value: dctools.MinimiseEmojiString(badgeString),
		})
	}

	if member.BoostedSince.IsValid() {
		field := dctools.UnixTimestamp(member.BoostedSince.Time())

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
		Value:  dctools.UnixTimestamp(user.CreatedAt()),
		Inline: true,
	})

	if member.Joined.IsValid() {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:   "Joined Server",
			Value:  dctools.UnixTimestamp(member.Joined.Time()),
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
	if user.Avatar != "" {
		avatarUploaded, _ = httputil.ImgUploadTime(user.AvatarURL())
	}
	if !avatarUploaded.IsZero() {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:   "Avatar Uploaded",
			Value:  dctools.UnixTimestamp(avatarUploaded),
			Inline: true,
		})
	}

	embed.Fields = append(embed.Fields, discord.EmbedField{
		Name:   "User ID",
		Value:  user.ID.String(),
		Inline: true,
	})

	if member.Banner != "" || member.User.Banner != "" {
		url := dctools.MemberBannerURL(member, guildID)
		embed.Image = &discord.EmbedImage{URL: url}
	}

	return &embed
}
