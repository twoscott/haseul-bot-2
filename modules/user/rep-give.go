package user

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var repGiveCommand = &router.SubCommand{
	Name:        "give",
	Description: "Gives a rep to a user",
	Handler: &router.CommandHandler{
		Executor: repGiveExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.UserOption{
			OptionName:  "user",
			Description: "The user to give a rep to",
			Required:    true,
		},
	},
}

func repGiveExec(ctx router.CommandCtx) {
	snowflake, _ := ctx.Options.Find("user").SnowflakeValue()
	if !snowflake.IsValid() {
		ctx.RespondError("Malformed snowflake provided.")
		return
	}

	senderID := ctx.Interaction.SenderID()
	targetID := discord.UserID(snowflake)

	if senderID == targetID {
		ctx.RespondWarning("You cannot rep yourself!")
		return
	}

	lastRep, err := db.Reps.GetUserLastRepTime(senderID, targetID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		ctx.RespondError("Error occurred while checking your recent reps.")
		return
	}

	cutoff := getRepCutoff()
	if !lastRep.Before(cutoff) {
		ctx.RespondWarning(
			"You cannot rep the same user more than once in the same day!",
		)
		return
	}

	remaining, err := db.Reps.GetUserRepsRemaining(ctx.Interaction.SenderID())
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while checking remaining reps.")
		return
	}

	if remaining == 0 {
		nextRepTime := getNextRepResetFromNow()
		ctx.RespondWarning(
			fmt.Sprintf(
				"You have no reps remaining! Your reps will be replenished %s.",
				dctools.UnixTimestampStyled(nextRepTime, dctools.RelativeTime),
			),
		)
		return
	}

	rep, err := db.Reps.RepUser(senderID, targetID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while attempting to rep user.")
		return
	}

	streak, err := db.Reps.GetUserStreak(senderID, targetID)
	if err != nil {
		log.Println(err)
	}

	target := ctx.Command.Resolved.Users[targetID]

	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: target.DisplayOrUsername(),
			Icon: target.AvatarURL(),
		},
		Fields: []discord.EmbedField{
			{Name: "Rep", Value: humanize.Comma(int64(rep)), Inline: true},
		},
		Color: dctools.EmbedBackColour,
	}

	days := streak.Days()
	if days > 0 {
		commaDays := humanize.Comma(int64(days))
		emojis := getStreakEmojiString(streak)

		hue := rand.Float64() * 360
		colour := dctools.HSVToColour(hue, 0.5, 0.9)
		embed.Color = colour

		embed.Fields = append(embed.Fields,
			discord.EmbedField{
				Name:   "Streak",
				Value:  fmt.Sprintf("%s days %s", commaDays, emojis),
				Inline: true,
			},
		)
	}

	message := fmt.Sprintf("You gave a rep to %s!", targetID.Mention())

	ctx.RespondSimple(message, embed)
}
