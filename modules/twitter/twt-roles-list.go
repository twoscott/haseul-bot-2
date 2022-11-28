package twitter

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var twtRolesListCommand = &router.SubCommand{
	Name:        "list",
	Description: "Lists all mention roles for a Twitter feed",
	Handler: &router.CommandHandler{
		Executor:      twtRoleListExec,
		Autocompleter: dbTwitterCompleter,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:   "twitter",
			Description:  "The Twitter user of the target feed",
			Required:     true,
			Autocomplete: true,
		},
		&discord.ChannelOption{
			OptionName:  "channel",
			Description: "The channel of the target feed",
			Required:    true,
			ChannelTypes: []discord.ChannelType{
				discord.GuildText,
				discord.GuildNews,
			},
		},
	},
}

func twtRoleListExec(ctx router.CommandCtx) {
	screenName := ctx.Options.Find("twitter").String()
	user, cerr := fetchUser(screenName)
	if cerr != nil {
		ctx.RespondCmdMessage(cerr)
		return
	}

	snowflake, _ := ctx.Options.Find("channel").SnowflakeValue()
	channelID := discord.ChannelID(snowflake)
	if !channelID.IsValid() {
		ctx.RespondWarning(
			"Malformed Discord channel provided.",
		)
		return
	}

	channel, err := ctx.State.Channel(channelID)
	if err != nil {
		ctx.RespondWarning("Invalid Discord channel provided.")
		return
	}

	_, err = db.Twitter.GetFeed(channel.ID, user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondWarning(
			fmt.Sprintf(
				"%s is not set up to receive tweets from @%s.",
				channel.Mention(), user.ScreenName,
			),
		)
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while fetching Twitter feed from the database.",
		)
		return
	}

	mentions, err := db.Twitter.GetMentions(channel.ID, user.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while fetching mention roles from the database.",
		)
		return
	}
	if len(mentions) < 1 {
		ctx.RespondWarning(
			"This Twitter feed has no mention roles set up.",
		)
		return
	}

	mentionList := make([]string, 0, len(mentions))
	for _, mention := range mentions {
		mentionList = append(mentionList, mention.RoleID.Mention())
	}

	descriptionPages := util.PagedLines(mentionList, 2048, 20)
	pages := make([]router.MessagePage, len(descriptionPages))
	footer := util.PluraliseWithCount("Mention Role", int64(len(mentions)))

	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		pages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Name: fmt.Sprintf(
							"@%s Twitter Mentions", user.ScreenName,
						),
						Icon: twitterIcon,
					},
					Description: description,
					Color:       twitterColour,
					Footer: &discord.EmbedFooter{
						Text: dctools.SeparateEmbedFooter(
							pageID,
							footer,
						),
					},
				},
			},
		}
	}

	ctx.RespondPaging(pages)
}
