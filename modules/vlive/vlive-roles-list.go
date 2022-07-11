package vlive

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

var vliveRolesListCommand = &router.SubCommand{
	Name:        "list",
	Description: "Lists all mention roles for a VLIVE feed",
	Handler: &router.CommandHandler{
		Executor:      vliveRoleListExec,
		Autocompleter: vliveFeedRemoveCompleter,
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
			Description: "The channel of the target feed",
			Required:    true,
			ChannelTypes: []discord.ChannelType{
				discord.GuildText,
				discord.GuildNews,
			},
		},
	},
}

func vliveRoleListExec(ctx router.CommandCtx) {
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

	channel, err := ctx.State.Channel(channelID)
	if err != nil {
		ctx.RespondWarning("Invalid Discord channel provided.")
		return
	}

	_, err = db.VLIVE.GetFeed(channel.ID, board.ID)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondWarning(
			fmt.Sprintf(
				"%s is not set up to receive VLIVE posts from %s.",
				channel.Mention(), vChannel.Name,
			),
		)
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while fetching VLIVE feed from the database.",
		)
		return
	}

	mentions, err := db.VLIVE.GetMentions(channel.ID, board.ID)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while fetching mention roles from the database.",
		)
		return
	}
	if len(mentions) < 1 {
		ctx.RespondWarning(
			"This VLIVE feed has no mention roles set up.",
		)
		return
	}

	mentionList := make([]string, 0, len(mentions))
	for _, mention := range mentions {
		mentionList = append(mentionList, mention.RoleID.Mention())
	}

	descriptionPages := util.PagedLines(mentionList, 2048, 20)
	messagePages := make([]router.MessagePage, len(descriptionPages))
	numOfMentions := len(mentions)
	numOfMentionsFooter := fmt.Sprintf(
		"%d %s",
		numOfMentions,
		util.Pluralise("Mention Role", int64(numOfMentions)),
	)

	for i, description := range descriptionPages {
		pageID := fmt.Sprintf("Page %d/%d", i+1, len(descriptionPages))
		messagePages[i] = router.MessagePage{
			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Name: fmt.Sprintf(
							"%s VLIVE Mentions", vChannel.Name,
						),
						Icon: vliveIcon,
					},
					Description: description,
					Color:       vliveColour,
					Footer: &discord.EmbedFooter{
						Text: dctools.SeparateEmbedFooter(
							pageID,
							numOfMentionsFooter,
						),
					},
				},
			},
		}
	}

	ctx.RespondPaging(messagePages)
}
