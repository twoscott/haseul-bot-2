package reminders

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/database/reminderdb"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"golang.org/x/exp/slices"
)

var remindersDeleteCommand = &router.SubCommand{
	Name:        "delete",
	Description: "Delete a reminder you previously set",
	Handler: &router.CommandHandler{
		Executor:      remindersDeleteExec,
		Autocompleter: completeReminderDelete,
		Ephemeral:     true,
	},
	Options: []discord.CommandOptionValue{
		&discord.IntegerOption{
			OptionName:   "reminder",
			Description:  "The reminder to delete",
			Required:     true,
			Autocomplete: true,
		},
	},
}

func remindersDeleteExec(ctx router.CommandCtx) {
	reminderID, _ := ctx.Options.Find("reminder").IntValue()

	ok, err := db.Reminders.DeleteForUser(ctx.Interaction.SenderID(), int32(reminderID))
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while trying to delete reminder.")
		return
	}
	if !ok {
		ctx.RespondWarning("I could not find this reminder.")
		return
	}

	ctx.RespondSuccess("Reminder deleted.")
}

func completeReminderDelete(ctx router.AutocompleteCtx) {
	reminders, err := db.Reminders.GetAllByUser(ctx.Interaction.SenderID())
	if err != nil {
		log.Println(err)
	}

	slices.SortFunc(reminders, func(a, b reminderdb.Reminder) bool {
		return a.Time.Unix() < b.Time.Unix()
	})

	choices := make(api.AutocompleteIntegerChoices, 0, len(reminders))
	for _, r := range reminders {
		choice := discord.IntegerChoice{
			Name:  fmt.Sprintln(dctools.EmbedTime(r.Time), "-", r.Content),
			Value: int(r.ID),
		}
		choices = append(choices, choice)
	}

	ctx.RespondChoices(choices)
}
