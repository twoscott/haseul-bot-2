package twitter

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var twtRetweetsToggleCommand = &router.SubCommand{
	Name:        "toggle",
	Description: "Toggles whether retweets are posted to a Twitter feed",
	Handler: &router.CommandHandler{
		Executor:      twtRetweetsToggleExec,
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

func twtRetweetsToggleExec(ctx router.CommandCtx) {
	screenName := ctx.Options.Find("user").String()
	user, cerr := fetchUser(screenName)
	if cerr != nil {
		ctx.RespondCmdMessage(cerr)
		return
	}

	snowflake, _ := ctx.Options.Find("channel").SnowflakeValue()
	channelID := discord.ChannelID(snowflake)
	if !channelID.IsValid() {
		ctx.RespondWarning("Malformed Discord channel provided.")
		return
	}

	channel, err := ctx.State.Channel(channelID)
	if err != nil {
		ctx.RespondWarning("Invalid Discord channel provided.")
		return
	}

	toggled, err := db.Twitter.ToggleFeedRetweets(channel.ID, user.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while toggling retweets.",
		)
		return
	}
	if !toggled {
		ctx.RespondWarning(
			fmt.Sprintf(
				"I could not find a Twitter feed for @%s in %s.",
				user.ScreenName, channel.Mention(),
			),
		)
		return
	}

	feed, err := db.Twitter.GetFeed(channel.ID, user.ID)
	if err != nil {
		ctx.RespondSuccess(
			"Retweets were toggled for this feed.",
		)
	}

	var content string
	if feed.Retweets {
		content = fmt.Sprintf(
			"You will now receive retweets from @%s in %s.",
			user.ScreenName, channel.Mention(),
		)
	} else {
		content = fmt.Sprintf(
			"You will no longer receive retweets from @%s in %s.",
			user.ScreenName, channel.Mention(),
		)
	}

	ctx.RespondSuccess(content)
}
