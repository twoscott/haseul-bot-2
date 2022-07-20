package vlive

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
)

var vliveFeedsRemoveCommand = &router.SubCommand{
	Name:        "remove",
	Description: "Removes a VLIVE feed from a Discord channel",
	Handler: &router.CommandHandler{
		Executor:      vliveFeedRemoveExec,
		Autocompleter: vliveFeedRemoveCompleter,
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
			Description:  "The VLIVE board to stop receiving VLIVE posts from",
			Required:     true,
			Autocomplete: true,
		},
		&discord.ChannelOption{
			OptionName:  "channel",
			Description: "The channel to stop receiving VLIVE posts in",
			Required:    true,
			ChannelTypes: []discord.ChannelType{
				discord.GuildText,
				discord.GuildNews,
			},
		},
	},
}

func vliveFeedRemoveExec(ctx router.CommandCtx) {
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

	channel, cerr := ctx.ParseAccessibleChannel(channelID)
	if cerr != nil {
		ctx.RespondCmdMessage(cerr)
		return
	}

	ok, err := db.VLIVE.RemoveFeed(channel.ID, board.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			fmt.Sprintf(
				"Error occurred while removing %s from the database.",
				vChannel.Name,
			),
		)
		return
	}
	if !ok {
		ctx.RespondWarning(
			fmt.Sprintf(
				"%s is not set up to receive VLIVE posts from %s.",
				channel.Mention(),
				vChannel.Name,
			),
		)
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf(
			"You will no longer receive VLIVE posts from %s's %s board in %s",
			vChannel.Name,
			board.Title,
			channel.Mention(),
		),
	)
}
