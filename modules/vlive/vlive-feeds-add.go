package vlive

import (
	"fmt"
	"log"
	"net/http"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/config"
	"github.com/twoscott/haseul-bot-2/database/vlivedb"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/cmdutil"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/vliveutil"
)

var vliveFeedsAddCommand = &router.SubCommand{
	Name:        "add",
	Description: "Adds a VLIVE feed to a Discord channel",
	Handler: &router.CommandHandler{
		Executor:      vliveFeedAddExec,
		Autocompleter: vliveFeedAddCompleter,
		Defer:         true,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:   "vlive-channel",
			Description:  "The VLIVE channel to search for VLIVE boards",
			Required:     true,
			Autocomplete: true,
		},
		&discord.IntegerOption{
			OptionName:   "vlive-board",
			Description:  "The VLIVE board to receive VLIVE posts from",
			Required:     true,
			Autocomplete: true,
		},
		&discord.ChannelOption{
			OptionName:  "channel",
			Description: "The channel to post VLIVE posts into",
			Required:    true,
			ChannelTypes: []discord.ChannelType{
				discord.GuildText,
				discord.GuildNews,
			},
		},
		&discord.IntegerOption{
			OptionName:  "types",
			Description: "The types of VLIVE posts to receive",
			Choices: []discord.IntegerChoice{
				{Name: "All Posts", Value: int(vlivedb.AllPosts)},
				{Name: "Videos Only", Value: int(vlivedb.VideosOnly)},
				{Name: "Posts Only", Value: int(vlivedb.PostsOnly)},
			},
		},
		&discord.BooleanOption{
			OptionName:  "reposts",
			Description: "Whether or not to receive reposts from the board",
		},
	},
}

func vliveFeedAddExec(ctx router.CommandCtx) {
	channelCode := ctx.Options.Find("vlive-channel").String()
	vChannel, cerr := fetchChannel(channelCode)
	if cerr != nil {
		ctx.RespondCmdMessage(cerr)
		return
	}

	boardID, err := ctx.Options.Find("vlive-board").IntValue()
	if err != nil {
		ctx.RespondWarning("Provided VLIVE board ID must be a number")
		return
	}
	board, cerr := fetchBoard(boardID)
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

	channel, cerr := cmdutil.ParseSendableChannel(ctx, channelID)
	if cerr != nil {
		ctx.RespondCmdMessage(cerr)
		return
	}

	postTypesInt, _ := ctx.Options.Find("types").IntValue()
	postTypes := vlivedb.PostTypes(postTypesInt)

	reposts, err := ctx.Options.Find("reposts").BoolValue()

	_, err = db.VLIVE.GetBoard(board.ID)
	if err != nil {
		cerr = addBoard(ctx, board)
	}
	if cerr != nil {
		ctx.RespondCmdMessage(cerr)
		return
	}

	cerr = checkGuildVLIVECount(&ctx, board.ID)
	if cerr != nil {
		ctx.RespondCmdMessage(cerr)
		return
	}

	ok, err := db.VLIVE.AddFeed(
		board.ID, ctx.Interaction.GuildID, channel.ID, postTypes, reposts,
	)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			fmt.Sprintf(
				"Error occurred while adding %s to the database.",
				vChannel.Name,
			),
		)
		return
	}
	if !ok {
		ctx.RespondWarning(
			fmt.Sprintf(
				"%s is already set up to receive VLIVE posts from %s.",
				channel.Mention(),
				vChannel.Name,
			),
		)
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf(
			"You will now receive VLIVE posts from %s in %s",
			vChannel.Name,
			channel.Mention(),
		),
	)
}

func addBoard(ctx router.CommandCtx, board *vliveutil.Board) router.CmdResponse {
	posts, res, err := vliveutil.GetBoardPosts(board.ID, 1)
	if err != nil {
		return router.Error("Error occured fetching VLIVE board posts.")
	}
	if res.StatusCode != http.StatusOK {
		return router.Errorf("Unable to fetch posts for %s.", board.Title)
	}

	pager := vliveutil.BoardPostsPager{}
	if len(posts) > 0 {
		post := posts[0]
		pager.PostID = post.ID
		pager.PostTimestamp = post.CreatedAt
	}

	_, err = db.VLIVE.AddBoard(
		board.ID, board.ChannelCode, pager.PostTimestamp, pager.PostID,
	)
	if err != nil {
		log.Println(err)
		return router.Errorf(
			"Error occurred while adding %s to the database.",
			board.Title,
		)
	}

	return nil
}

func checkGuildVLIVECount(
	ctx *router.CommandCtx, boardID int64) router.CmdResponse {

	cfg := config.GetInstance()
	if ctx.Interaction.GuildID == cfg.Bot.RootGuildID {
		return nil
	}

	vliveCount, err := db.VLIVE.GetGuildBoardCount(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		return router.Error(
			"Error occurred while checking current VLIVE feeds.",
		)
	}

	if vliveCount < 1 {
		return nil
	}

	guild, err := ctx.State.Guild(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		return router.Error("Error occurred while fetching server.")
	}

	patron, err := pat.GetActiveDiscordPatron(guild.OwnerID)
	if err != nil {
		log.Println(err)
	}

	if patron == nil || patron.CurrentlyEntitledAmountCents < 300 {
		if vliveCount >= 3 {
			return router.Warning(
				"You can only add VLIVE feeds for 3 VLIVE boards per server.",
			)
		}
	} else if patron != nil && patron.CurrentlyEntitledAmountCents >= 300 {
		if vliveCount >= 10 {
			return router.Warning(
				"You can only add VLIVE feeds for 10 VLIVE boards per server.",
			)
		}
	}

	return nil
}

func vliveFeedAddCompleter(ctx router.AutocompleteCtx) {
	switch ctx.Focused.Name {
	case "vlive-channel":
		completeVliveChannel(ctx)
	case "vlive-board":
		completeVliveBoard(ctx)
	}
}

func completeVliveChannel(ctx router.AutocompleteCtx) {
	query := ctx.Focused.String()

	channels, res, err := vliveutil.SearchChannels(query, 10)
	if err != nil {
		log.Println(err)
		return
	}
	if res.StatusCode != http.StatusOK {
		return
	}

	choices := make(api.AutocompleteStringChoices, 0, len(channels))
	for _, c := range channels {
		choices = append(choices, discord.StringChoice{
			Name: c.Name, Value: c.Code,
		})
	}

	ctx.RespondChoices(choices)
}

func completeVliveBoard(ctx router.AutocompleteCtx) {
	channelCode := ctx.Options.Find("vlive-channel").String()
	query := ctx.Focused.String()

	boards, res, err := vliveutil.GetUnwrappedBoards(channelCode)
	if err != nil {
		log.Println(err)
		return
	}
	if res.StatusCode != http.StatusOK {
		return
	}

	choices := make(api.AutocompleteIntegerChoices, 0, len(boards))
	for _, b := range boards {
		choices = append(choices, discord.IntegerChoice{
			Name: b.Title, Value: int(b.ID),
		})
	}

	choices = dctools.SearchSortIntChoices(choices, query)
	ctx.RespondChoices(choices)
}
