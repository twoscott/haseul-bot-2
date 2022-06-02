package youtube

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/ytutil"
)

var youTubeSearchCommand = &router.SubCommand{
	Name:        "search",
	Description: "Searches YouTube for a query and displays results",
	Handler: &router.CommandHandler{
		Executor:      youTubeSearchExec,
		Autocompleter: youTubeSearchCompleter,
		Defer:         true,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:   "query",
			Description:  "The query to search YouTube for",
			Required:     true,
			Autocomplete: true,
		},
		&discord.IntegerOption{
			OptionName:  "results",
			Description: "The amount of results to return",
			Min:         option.NewInt(1),
			Max:         option.NewInt(20),
		},
	},
}

func youTubeSearchExec(ctx router.CommandCtx) {
	query := ctx.Options.Find("query").String()
	results, _ := ctx.Options.Find("results").IntValue()
	if results == 0 {
		results = 20
	}

	videoLinks, err := ytutil.MultiSearch(query, int(results))
	if err == ytutil.ErrNoResultsFound {
		ctx.RespondWarning(
			fmt.Sprintf("No results found for '%s'.", query),
		)
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred trying to fetch YouTube results.",
		)
		return
	}

	err = db.YouTube.AddHistoryAndClear(
		ctx.Interaction.SenderID(), ctx.Interaction.ID, query,
	)
	if err != nil {
		log.Println(err)
	}

	messagePages := make([]router.MessagePage, len(videoLinks))
	for i, link := range videoLinks {
		messagePages[i] = router.MessagePage{
			Content: fmt.Sprintf("%d. %s", i+1, link),
		}
	}

	ctx.RespondConfirmationPaging(messagePages)
}

func youTubeSearchCompleter(ctx router.AutocompleteCtx) {
	query := ctx.Focused.String()

	var (
		suggestions []string
		err         error
	)

	if query == "" {
		suggestions, err = db.YouTube.GetHistory(ctx.Interaction.SenderID())
	} else {
		suggestions, err = ytutil.GetSuggestions(query)
	}
	if err != nil {
		log.Println(err)
		return
	}
	if len(suggestions) < 1 {
		suggestions = append(suggestions, query)
	}

	choices := make([]api.AutocompleteChoice, 0, len(suggestions))
	for _, s := range suggestions {
		choice := api.AutocompleteChoice{Name: s, Value: s}
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
