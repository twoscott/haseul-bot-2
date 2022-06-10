package twitter

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/config"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var twtFeedsAddCommand = &router.SubCommand{
	Name:        "add",
	Description: "Adds a Twitter feed to a Discord channel",
	Handler: &router.CommandHandler{
		Executor:      twtFeedAddExec,
		Autocompleter: twtFeedAddCompleter,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:   "twitter",
			Description:  "The Twitter user to listen for Tweets from",
			Required:     true,
			Autocomplete: true,
		},
		&discord.ChannelOption{
			OptionName:  "channel",
			Description: "The channel to post Tweets from the user into",
			Required:    true,
			ChannelTypes: []discord.ChannelType{
				discord.GuildText,
				discord.GuildNews,
			},
		},
	},
}

func twtFeedAddExec(ctx router.CommandCtx) {
	screenName := ctx.Options.Find("user").String()
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

	channel, cerr := parseChannelArg(ctx, channelID)
	if cerr != nil {
		ctx.RespondCmdMessage(cerr)
		return
	}

	botUser, err := ctx.State.Me()
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			fmt.Sprintf("Error occurred checking my permissions in %s.",
				channel.Mention(),
			),
		)
		return
	}

	botPermissions, err := ctx.State.Permissions(channel.ID, botUser.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			fmt.Sprintf("Error occurred checking my permissions in %s.",
				channel.Mention(),
			),
		)
		return
	}

	neededPerms := dctools.PermissionsBitfield(
		discord.PermissionViewChannel,
		discord.PermissionSendMessages,
	)

	if !botPermissions.Has(neededPerms) {
		ctx.RespondWarning(
			fmt.Sprintf("I do not have permission to send messages in %s!",
				channel.Mention(),
			),
		)
		return
	}

	_, err = db.Twitter.GetUser(user.ID)
	if err != nil {
		cerr = addUser(ctx, user)
	}
	if cerr != nil {
		ctx.RespondCmdMessage(cerr)
		return
	}

	_, err = db.Twitter.GetUserByGuild(channel.GuildID, user.ID)
	if err != nil {
		cerr = checkGuildTwitterCount(&ctx, user.ID)
	}
	if cerr != nil {
		ctx.RespondCmdMessage(cerr)
		return
	}

	ok, err := db.Twitter.AddFeed(ctx.Interaction.GuildID, channel.ID, user.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			fmt.Sprintf(
				"Error occurred while adding @%s to the database.",
				user.ScreenName,
			),
		)
		return
	}
	if !ok {
		ctx.RespondWarning(
			fmt.Sprintf(
				"%s is already set up to receive tweets from @%s.",
				channel.Mention(), user.ScreenName,
			),
		)
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf(
			"You will now receive tweets from @%s in %s.",
			user.ScreenName, channel.Mention()),
	)
}

func addUser(ctx router.CommandCtx, user *twitter.User) router.CmdResponse {
	tweets, resp, err := twt.Timelines.UserTimeline(&twitter.UserTimelineParams{
		UserID:          user.ID,
		Count:           1,
		ExcludeReplies:  twitter.Bool(false),
		IncludeRetweets: twitter.Bool(true),
		TrimUser:        twitter.Bool(true),
	})
	if err != nil {
		return router.Error(
			"Unknown error occurred while trying to fetch tweets.",
		)
	}
	if resp.StatusCode != http.StatusOK {
		log.Println(err)
		return router.Errorf(
			"Error occurred while fetching neccesary data from @%s.",
			user.ScreenName,
		)
	}

	var lastTweetID int64
	if len(tweets) < 1 {
		lastTweetID = 0
	} else {
		lastTweetID = tweets[0].ID
	}

	_, err = db.Twitter.AddUser(user.ID, lastTweetID)
	if err != nil {
		log.Println(err)
		return router.Errorf(
			"Error occurred while adding @%s to the database.",
			user.ScreenName,
		)
	}

	return nil
}

func checkGuildTwitterCount(
	ctx *router.CommandCtx, twitterUserID int64) router.CmdResponse {

	cfg := config.GetInstance()
	if ctx.Interaction.GuildID == cfg.Bot.RootGuildID {
		return nil
	}

	twitterCount, err := db.Twitter.GetGuildUserCount(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		return router.Error(
			"Error occurred while checking current Twitter feeds.",
		)
	}

	if twitterCount < 1 {
		return nil
	}

	guild, err := ctx.State.GuildWithCount(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		return router.Error("Error occurred while checking member count.")
	}

	memberCount := guild.ApproximateMembers
	patron, err := pat.GetActiveDiscordPatron(guild.OwnerID)
	if err != nil {
		log.Println(err)
		return router.Error(
			"Error occurred while checking server owner's patron status.",
		)
	}

	if patron == nil || patron.CurrentlyEntitledAmountCents < 300 {
		if memberCount < 100 && twitterCount > 0 {
			return router.Error(
				"Your server must have at least 100 members " +
					"to add feeds for more than 1 Twitter account.",
			)
		} else if memberCount < 250 && twitterCount > 1 {
			return router.Error(
				"Your server must have at least 250 members " +
					"to add feeds for more than 2 Twitter accounts.",
			)
		}
	}
	if patron == nil || patron.CurrentlyEntitledAmountCents < 1000 {
		if twitterCount > 2 {
			return router.Error(
				"Your server cannot have feeds for more than " +
					"3 Twitter accounts at once.",
			)
		}
	}
	if patron != nil && patron.CurrentlyEntitledAmountCents >= 1000 {
		if twitterCount > 9 {
			return router.Error(
				"Your server cannot have feeds for more than " +
					"10 Twitter accounts at once.",
			)
		}
	}

	return nil
}

func twtFeedAddCompleter(ctx router.AutocompleteCtx) {
	user := ctx.Options.Find("twitter").String()
	if user == "" {
		return
	}

	users, resp, err := twt.Users.Search(user, &twitter.UserSearchParams{
		Query:           user,
		Page:            1,
		Count:           10,
		IncludeEntities: twitter.Bool(false),
	})
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}

	choices := make(api.AutocompleteStringChoices, 0, len(users))
	for _, u := range users {
		choice := discord.StringChoice{
			Name: u.ScreenName, Value: u.ScreenName,
		}
		choices = append(choices, choice)
	}

	ctx.State.RespondInteraction(ctx.Interaction.ID, ctx.Interaction.Token,
		api.InteractionResponse{
			Type: api.AutocompleteResult,
			Data: &api.InteractionResponseData{
				Choices: &choices,
			},
		},
	)
}
