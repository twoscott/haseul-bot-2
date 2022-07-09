package vlive

import (
	"log"
	"net/http"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/vliveutil"
)

const (
	vliveIcon   = "https://www.vlive.tv/favicon.ico?2020091516"
	vliveColour = 0x54f7ff
)

func fetchChannel(channelCode string) (*vliveutil.Channel, router.CmdResponse) {
	channel, res, err := vliveutil.GetChannel(channelCode)
	if err != nil {
		return nil, router.Error("Error occurred fetching VLIVE channel.")
	}

	switch res.StatusCode {
	case http.StatusOK:
		return channel, nil
	case http.StatusNotFound:
		return nil, router.Warning("This VLIVE channel does not exist.")
	default:
		return nil, router.Error("Unable to fetch VLIVE channel.")
	}
}

func fetchBoard(boardID int64) (*vliveutil.Board, router.CmdResponse) {
	board, res, err := vliveutil.GetBoard(boardID)
	if err != nil {
		return nil, router.Error("Error occurred fetching VLIVE board.")
	}

	switch res.StatusCode {
	case http.StatusOK:
		return board, nil
	case http.StatusNotFound:
		return nil, router.Warning("This VLIVE board does not exist.")
	default:
		return nil, router.Error("Unable to fetch VLIVE board.")
	}
}

func vliveFeedRemoveCompleter(ctx router.AutocompleteCtx) {
	switch ctx.Focused.Name {
	case "vlive-channel":
		completeVliveChannelDB(ctx)
	case "vlive-board":
		completeVliveBoardDB(ctx)
	}
}

func completeVliveChannelDB(ctx router.AutocompleteCtx) {
	query := ctx.Focused.String()

	channelCodes, err := db.VLIVE.GetChannelCodesByGuild(
		ctx.Interaction.GuildID,
	)
	if err != nil {
		log.Println(err)
		return
	}

	channels := make([]*vliveutil.Channel, 0, len(channelCodes))
	for _, code := range channelCodes {
		channel, res, err := vliveutil.GetChannel(code)
		if err != nil {
			log.Println(err)
			continue
		}
		if res.StatusCode != http.StatusOK {
			continue
		}

		channels = append(channels, channel)
	}

	choices := make(api.AutocompleteStringChoices, 0, len(channels))
	for _, c := range channels {
		choices = append(choices, discord.StringChoice{
			Name: c.Name, Value: c.Code,
		})
	}

	choices = dctools.SearchSortStringChoices(choices, query)
	ctx.RespondChoices(choices)
}

func completeVliveBoardDB(ctx router.AutocompleteCtx) {
	channelCode := ctx.Options.Find("vlive-channel").String()
	query := ctx.Focused.String()

	dbBoards, err := db.VLIVE.GetBoardsByVLIVEChannel(channelCode)
	if err != nil {
		log.Println(err)
		return
	}

	boards := make([]*vliveutil.Board, 0, len(dbBoards))
	for _, b := range dbBoards {
		board, res, err := vliveutil.GetBoard(b.ID)
		if err != nil {
			log.Println(err)
			continue
		}
		if res.StatusCode != http.StatusOK {
			continue
		}

		boards = append(boards, board)
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
