package router

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
)

// AutocompleteCtx wraps router and includes data about the autocomplete
// interaction to be passed to the receiving completion handler.
type AutocompleteCtx struct {
	*InteractionCtx
	Options discord.AutocompleteOptions
	// Focused is the option that is currently being typed in by the user.
	Focused discord.AutocompleteOption
}

// AutocompleteInteractionKey returns the string representing the command and
// subcommands of an autocomplete interaction as a single string,
// for use in the commandHandlers hash map.
func AutocompleteInteractionKey(
	completion *discord.AutocompleteInteraction) string {

	if len(completion.Options) < 1 {
		return completion.Name
	}

	return completion.Name + autocompleteString(&completion.Options[0])
}

func autocompleteString(option *discord.AutocompleteOption) string {
	switch option.Type {
	case discord.SubcommandGroupOptionType:
		return "/" + option.Name + autocompleteString(&option.Options[0])
	case discord.SubcommandOptionType:
		return "/" + option.Name
	default:
		return ""
	}
}

func (ctx AutocompleteCtx) RespondChoices(
	choices api.AutocompleteChoices) error {

	return ctx.State.RespondInteraction(
		ctx.Interaction.ID,
		ctx.Interaction.Token,
		api.InteractionResponse{
			Type: api.AutocompleteResult,
			Data: &api.InteractionResponseData{
				Choices: choices,
			},
		},
	)
}
