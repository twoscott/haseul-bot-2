package user

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

var divorceButtons = &discord.ActionRowComponent{
	&discord.ButtonComponent{
		Label:    "No",
		CustomID: "NO",
		Emoji: &discord.ComponentEmoji{
			Name: "💗",
		},
		Style: discord.SecondaryButtonStyle(),
	},
	&discord.ButtonComponent{
		Label:    "Yes",
		CustomID: "YES",
		Emoji: &discord.ComponentEmoji{
			Name: "💔",
		},
		Style: discord.PrimaryButtonStyle(),
	},
}

var marriageDivorceCommand = &router.SubCommand{
	Name:        "divorce",
	Description: "Divorce your spouse",
	Handler: &router.CommandHandler{
		Executor: marriageDivorceExec,
	},
}

func marriageDivorceExec(ctx router.CommandCtx) {
	marriage, err := db.Marriages.GetUserMarriage(ctx.Interaction.SenderID())
	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondWarning("You are not married to anyone!")
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

	msg := fmt.Sprintf("Are you sure you want to divorce %s?", spouseName)

	err = ctx.RespondMessage(api.InteractionResponseData{
		Content:    option.NewNullableString(msg),
		Components: discord.ComponentsPtr(divorceButtons),
	})
	if err != nil {
		ctx.RespondError("Error occurred while divorcing.")
		return
	}

	ch, cancel := ctx.State.ChanFor(func(ev any) bool {
		if i, ok := ev.(*gateway.InteractionCreateEvent); ok {
			if _, ok := i.Data.(*discord.ButtonInteraction); ok {
				if i.Message.Interaction.ID == ctx.Interaction.ID && i.SenderID() == ctx.Interaction.SenderID() {
					return true
				}

				ctx.State.RespondInteraction(i.ID, i.Token,
					api.InteractionResponse{
						Type: api.DeferredMessageUpdate,
					},
				)
			}
		}
		return false
	})

	select {
	case ev := <-ch:
		if i, ok := ev.(*gateway.InteractionCreateEvent); ok {
			if _, ok := i.Data.(*discord.ButtonInteraction); ok {
				if i.Message.Interaction.ID == ctx.Interaction.ID && i.SenderID() == ctx.Interaction.SenderID() {
					itx := &router.InteractionCtx{
						Router:      ctx.Router,
						Interaction: &i.InteractionEvent,
					}

					if i.Data.(*discord.ButtonInteraction).CustomID == "YES" {
						acceptDivorce(itx, spouseName)
					}
					rejectDivorce(itx)
				}
			}
		}
		disableDivorceButtons(ctx)
		cancel()

	case <-time.After(24 * time.Hour):
		disableDivorceButtons(ctx)
		delete(proposals, ctx.Interaction.SenderID())
		cancel()
	}
}

func disableDivorceButtons(ctx router.CommandCtx) {
	d := dctools.DisabledButtons(*divorceButtons)
	ctx.State.EditInteractionResponse(
		ctx.Interaction.AppID,
		ctx.Interaction.Token,
		api.EditInteractionResponseData{
			Components: discord.ComponentsPtr(&d),
		},
	)
}

func acceptDivorce(itx *router.InteractionCtx, spouseName string) {
	removed, err := db.Marriages.Remove(itx.Interaction.SenderID())
	if err != nil {
		log.Println(err)
		itx.RespondError("Error occurred while divorcing.")
		return
	}
	if removed == 0 {
		itx.RespondWarning("You are not married to anyone!")
		return
	}

	itx.RespondTextf("You have divorced %s.", spouseName)
}

func rejectDivorce(itx *router.InteractionCtx) {
	itx.RespondText("You have cancelled the divorce.")
}
