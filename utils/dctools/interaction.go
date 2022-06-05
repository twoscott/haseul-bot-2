package dctools

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
)

// MessageResponse returns a message interaction with source response object
// containing the supplied data.
func MessageResponse(
	data api.InteractionResponseData) *api.InteractionResponse {
	return &api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &data,
	}
}

// UpdateMessageResponse returns an update message response object containing
// the supplied data.
func UpdateMessageResponse(
	data api.InteractionResponseData) *api.InteractionResponse {

	return &api.InteractionResponse{
		Type: api.UpdateMessage,
		Data: &data,
	}
}

// DeferredMessageResponse returns a deferred message interaction with source
// response object containing the supplied data.
func DeferredMessageResponse(
	data api.InteractionResponseData) *api.InteractionResponse {

	return &api.InteractionResponse{
		Type: api.DeferredMessageInteractionWithSource,
		Data: &data,
	}
}

// IsValueOption returns whether or not the option type is a value type or not
// (if the type is neither a sub command group or sub command).
func IsValueOption(optionType discord.CommandOptionType) bool {
	return optionType != discord.SubcommandOptionType &&
		optionType != discord.SubcommandGroupOptionType
}

// CommandOptions returns the options of the deepest command or subcommand in
// the chain containing only values, ignoring the need to make nested
// Find() calls to access options the user entered after the subcommand
// or command.
func CommandOptions(
	command *discord.CommandInteraction) discord.CommandInteractionOptions {

	return deepestCommandOptions(command.Options)
}

func deepestCommandOptions(
	opts discord.CommandInteractionOptions) discord.CommandInteractionOptions {

	if len(opts) < 1 || IsValueOption(opts[0].Type) {
		return opts
	}

	return deepestCommandOptions(opts[0].Options)
}

// AutocompleteOptions returns the options of the deepest command or subcommand
// in the chain containing only option values, ignoring the need to make nested
// loops to find the Focused option
func AutocompleteOptions(
	completion *discord.AutocompleteInteraction) discord.AutocompleteOptions {

	return deepestAutocompleteOptions(completion.Options)
}

func deepestAutocompleteOptions(
	opts []discord.AutocompleteOption) discord.AutocompleteOptions {

	if len(opts) < 1 || IsValueOption(opts[0].Type) {
		return opts
	}

	return deepestAutocompleteOptions(opts[0].Options)
}

// FocusedOption finds the currently focused option from the autocomplete
// interaction sent from Discord.
func FocusedOption(
	completion *discord.AutocompleteInteraction) *discord.AutocompleteOption {

	return findFocusedOption(completion.Options)
}

func findFocusedOption(
	options []discord.AutocompleteOption) *discord.AutocompleteOption {

	for _, opt := range options {
		switch opt.Type {
		case discord.SubcommandGroupOptionType, discord.SubcommandOptionType:
			return findFocusedOption(opt.Options)
		default:
			if opt.Focused {
				return &opt
			}
		}
	}

	return nil
}

// MakeStringChoices takes a slice of strings and turns them into a slice of
// autocomplete choices with matching names and string values.
func MakeStringChoices(choiceStrings []string) api.AutocompleteStringChoices {
	choices := make(api.AutocompleteStringChoices, 0, len(choiceStrings))
	for _, c := range choiceStrings {
		choice := discord.StringChoice{Name: c, Value: c}
		choices = append(choices, choice)
	}

	return choices
}
