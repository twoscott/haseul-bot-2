package reminders

import (
	"fmt"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

const reminderLimit = 25

var remindersAddCommand = &router.SubCommand{
	Name:        "add",
	Description: "Get reminded about something in the future",
	Handler: &router.CommandHandler{
		Executor:  remindersAddExec,
		Ephemeral: true,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:  "duration",
			Description: "How long from now to wait until reminding you",
			MaxLength:   option.NewInt(64),
			Required:    true,
		},
		&discord.StringOption{
			OptionName:  "reminder",
			Description: "What to be reminded of",
			MaxLength:   option.NewInt(2048),
			Required:    true,
		},
	},
}

func remindersAddExec(ctx router.CommandCtx) {
	pending, err := db.Reminders.GetAllByUser(ctx.Interaction.SenderID())
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while checking pending reminders.")
		return
	}

	if len(pending) >= reminderLimit {
		ctx.RespondWarning(
			fmt.Sprintf(
				"You cannot have more than %d pending reminders at once.",
				reminderLimit,
			),
		)
		return
	}

	startTime := time.Now()
	durationString := ctx.Options.Find("duration").String()
	reminder := ctx.Options.Find("reminder").String()

	timePeriod := util.ParseTimePeriod(durationString)
	if timePeriod.IsNull() {
		ctx.RespondWarning(
			"Invalid time period given. Example format: `3 days 4hr 6 min 2s`",
		)
		return
	}

	newTime := timePeriod.AfterTime(startTime)

	dmChannel, err := ctx.State.CreatePrivateChannel(ctx.Interaction.SenderID())
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while trying to DM you.",
		)
		return
	}

	reminderId, err := db.Reminders.Add(ctx.Interaction.SenderID(), newTime, reminder)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while adding reminder to the database.",
		)
		return
	}

	dmMsg := fmt.Sprintf(
		"You will be be reminded to '%s' on %s.",
		reminder, dctools.UnixTimestamp(newTime),
	)

	_, err = ctx.State.SendMessage(dmChannel.ID, dmMsg)
	if dctools.ErrCannotDM(err) {
		ctx.RespondWarning(
			"I am unable to DM you. " +
				"Please open your DMs to server members in your settings.",
		)
		db.Reminders.DeleteForUser(ctx.Interaction.SenderID(), reminderId)
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf(
			"Reminder set for %s.",
			dctools.UnixTimestampStyled(newTime, dctools.LongDateTime),
		),
	)
}
