package user

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var repGiveCommand = &router.SubCommand{
	Name:        "give",
	Description: "Gives a rep to a user",
	Handler: &router.CommandHandler{
		Executor: repGiveExec,
		Defer:    true,
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
				"You have no reps remaining! Your reps will be replenished %s",
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

	log.Println(streak)

	target := ctx.Command.Resolved.Users[targetID]

	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: target.Username,
			Icon: target.AvatarURL(),
		},
		Fields: []discord.EmbedField{
			{Name: "Rep", Value: strconv.Itoa(rep)},
		},
		Color: dctools.EmbedBackColour,
	}

	if streak > 0 {
		embed.Fields = append(
			embed.Fields,
			discord.EmbedField{Name: "Streak", Value: strconv.Itoa(streak)},
		)
	}

	message := fmt.Sprintf("You gave a rep to %s!", targetID.Mention())

	ctx.RespondSimple(message, embed)
}
