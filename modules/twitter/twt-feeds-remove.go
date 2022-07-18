package twitter

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
)

var twtFeedsRemoveCommand = &router.SubCommand{
	Name:        "remove",
	Description: "Removes a Twitter feed from a Discord channel",
	Handler: &router.CommandHandler{
		Executor:      twtFeedRemoveExec,
		Autocompleter: dbTwitterCompleter,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:   "twitter",
			Description:  "The Twitter user to stop receiving Tweets from",
			Required:     true,
			Autocomplete: true,
		},
		&discord.ChannelOption{
			OptionName:  "channel",
			Description: "The channel to stop posting Tweets into",
			Required:    true,
			ChannelTypes: []discord.ChannelType{
				discord.GuildText,
				discord.GuildNews,
			},
		},
	},
}

func twtFeedRemoveExec(ctx router.CommandCtx) {
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

	channel, cerr := cmdutil.ParseAccessibleChannel(ctx, channelID)
	if cerr != nil {
		ctx.RespondCmdMessage(cerr)
		return
	}

	_, err := db.Twitter.RemoveMentions(channel.ID, user.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while removing roles from the database.",
		)
		return
	}

	ok, err := db.Twitter.RemoveFeed(channel.ID, user.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			fmt.Sprintf(
				"Error occurred while removing @%s from the database.",
				user.ScreenName,
			),
		)
		return
	}
	if !ok {
		ctx.RespondWarning(
			fmt.Sprintf(
				"%s is not set up to receive tweets from @%s.",
				channel.Mention(), user.ScreenName,
			),
		)
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf(
			"You will no longer receive tweets from @%s in %s.",
			user.ScreenName, channel.Mention(),
		),
	)
}
