package twitter

import (
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var twtRoleListCommand = &router.Command{
	Name:                "list",
	RequiredPermissions: discord.PermissionManageChannels,
	IncludeAdmin:        true,
	UseTyping:           true,
	Run:                 twtRoleListRun,
}

func twtRoleListRun(ctx router.CommandCtx, _ []string) {
	dbUsers, err := db.Twitter.GetUsersByGuild(ctx.Msg.GuildID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while fetching Twitter users from the database.",
		)
		return
	}

	mentions, err := db.Twitter.GetMentionsByGuild(ctx.Msg.GuildID)
	if err != nil {
		log.Println(err)
		dctools.ReplyWithError(ctx.State, ctx.Msg,
			"Error occurred while fetching mention roles from the database.",
		)
		return
	}

	if len(mentions) < 1 {
		dctools.ReplyWithWarning(ctx.State, ctx.Msg,
			"This server has no Twitter mentions set up.",
		)
		return
	}

	screenNames := make(map[int64]string)
	for _, dbUser := range dbUsers {
		user, _, _ := twt.Users.Show(&twitter.UserShowParams{UserID: dbUser.ID})
		if user != nil {
			screenNames[dbUser.ID] = user.ScreenName
		}
	}

	mentionList := make([]string, 0, len(mentions))
	for _, mention := range mentions {
		mentionEntry := mention.RoleID.Mention() + " - "

		screenName := screenNames[mention.TwitterUserID]
		if screenName == "" {
			mentionEntry += fmt.Sprintf("(ID#%d)", mention.TwitterUserID)
		} else {
			handle := fmt.Sprintf("@%s", screenName)
			url := fmt.Sprintf("https://twitter.com/%s", screenName)
			mentionEntry += dctools.Hyperlink(handle, url)
		}

		mentionEntry += " " + mention.ChannelID.Mention()

		mentionList = append(mentionList, mentionEntry)
	}

	descriptionPages := util.PagedLines(mentionList, 2048, 20)
	messagePages := make([]router.MessagePage, len(descriptionPages))
	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		messagePages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Name: "Twitter Mentions", Icon: twitterIcon,
					},
					Description: description,
					Color:       twitterColour,
					Footer:      &discord.EmbedFooter{Text: pageID},
				},
			},
		}
	}

	cmdutil.ReplyWithPaging(ctx, ctx.Msg, messagePages)
}
