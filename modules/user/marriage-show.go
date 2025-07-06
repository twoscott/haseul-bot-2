package user

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var marriageShowCommand = &router.SubCommand{
	Name:        "show",
	Description: "Displays information about your marriage",
	Handler: &router.CommandHandler{
		Executor: marriageShowExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.UserOption{
			OptionName:  "user",
			Description: "The user to show marriage information for",
		},
	},
}

func marriageShowExec(ctx router.CommandCtx) {
	snowflake, _ := ctx.Options.Find("user").SnowflakeValue()
	userID := discord.UserID(snowflake)

	if !userID.IsValid() {
		userID = ctx.Interaction.SenderID()
	}

	marriage, err := db.Marriages.GetUserMarriage(userID)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondWarning("This user is not married to anyone!")
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching marriage data.")
		return
	}

	spouseID := marriage.Spouse(ctx.Interaction.SenderID())
	spouseName := spouseID.Mention()

	spouse, err := ctx.State.User(spouseID)
	if err == nil {
		spouseName = spouse.DisplayOrUsername()
	}

	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: fmt.Sprintf("%s and %s", ctx.Interaction.Sender().DisplayName, spouseName),
		},
		Description: fmt.Sprintf(
			"💒 %s and %s got married %s 💐",
			ctx.Interaction.Sender().DisplayName,
			spouseName,
			dctools.TimestampStyled(marriage.MarriedAt, dctools.RelativeTime),
		),
		Footer: &discord.EmbedFooter{
			Text: "💗 Married",
		},
		Timestamp: discord.Timestamp(marriage.MarriedAt),
		Color:     marriageColour,
	}

	ctx.RespondEmbed(embed)
}
