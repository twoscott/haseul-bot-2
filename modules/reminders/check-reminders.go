package reminders

import (
	"log"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/twoscott/haseul-bot-2/database/reminderdb"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

const interval = time.Second * 30

func startCheckingReminders(st *state.State) {
	for {
		start := time.Now()
		log.Println("Started checking reminders")

		checkReminders(st)

		elapsed := time.Since(start)
		log.Printf(
			"Finished checking reminders, took: %1.2fs\n", elapsed.Seconds(),
		)

		wait := interval - elapsed
		<-time.After(wait)
	}
}

func checkReminders(st *state.State) {
	reminders, err := db.Reminders.GetOverdueReminders()
	if err != nil {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup
	for _, reminder := range reminders {
		wg.Add(1)
		go func(r reminderdb.Reminder) {
			defer wg.Done()
			sendReminder(st, r)
		}(reminder)
	}
}

func sendReminder(st *state.State, reminder reminderdb.Reminder) {
	dmChannel, err := st.CreatePrivateChannel(reminder.UserID)
	if err != nil {
		log.Println(err)
		return
	}

	msg := "⏰ Reminder has been triggered."
	_, err = st.SendMessage(dmChannel.ID, msg, discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: "Reminder",
		},
		Description: reminder.Content,
		Color:       dctools.EmbedBackColour,
		Footer: &discord.EmbedFooter{
			Text: "Reminder set on",
		},
		Timestamp: discord.Timestamp(reminder.Created),
	})
	if err != nil {
		log.Println(err)
		return
	}

	db.Reminders.DeleteForUser(reminder.UserID, reminder.ID)
}
