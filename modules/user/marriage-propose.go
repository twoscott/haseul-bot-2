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
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var proposalButtons = &discord.ActionRowComponent{
	&discord.ButtonComponent{
		Label:    "No",
		CustomID: "NO",
		Emoji: &discord.ComponentEmoji{
			Name: "💔",
		},
		Style: discord.SecondaryButtonStyle(),
	},
	&discord.ButtonComponent{
		Label:    "Yes",
		CustomID: "YES",
		Emoji: &discord.ComponentEmoji{
			Name: "💍",
		},
		Style: discord.SuccessButtonStyle(),
	},
}

var proposals = make(map[discord.UserID]discord.UserID)

var marriageProposeCommand = &router.SubCommand{
	Name:        "propose",
	Description: "Propose to a user",
	Handler: &router.CommandHandler{
		Executor: marriageProposeExec,
	},
	Options: []discord.CommandOptionValue{
		&discord.UserOption{
			OptionName:  "user",
			Description: "The user to propose to",
			Required:    true,
		},
	},
}

func marriageProposeExec(ctx router.CommandCtx) {
	proposeeSnowflake, _ := ctx.Options.Find("user").SnowflakeValue()
	proposeeID := discord.UserID(proposeeSnowflake)

	if proposeeID == ctx.Interaction.SenderID() {
		ctx.RespondWarning("You cannot propose to yourself!")
		return
	}

	proposee := ctx.Command.Resolved.Users[proposeeID]
	proposeeAvatar := proposee.AvatarURL()

	marriage, err := db.Marriages.GetUserMarriage(ctx.Interaction.SenderID())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching marriage data.")
		return
	}
	if marriage.Spouse(ctx.Interaction.SenderID()) == proposeeID {
		ctx.RespondWarningf("You and %s are already married!", proposee.DisplayOrUsername())
		return
	}
	if s := marriage.Spouse(ctx.Interaction.SenderID()); s.IsValid() {
		spouse, err := ctx.State.User(s)
		if err != nil {
			ctx.RespondWarning("You are already married!")
			return
		}

		ctx.RespondWarningf(
			"You are already married to %s (@%s)!", spouse.DisplayName, spouse.Username,
		)
		return
	}
	if s := marriage.Spouse(proposeeID); s.IsValid() {
		spouse, err := ctx.State.User(s)
		if err != nil {
			ctx.RespondWarningf("%s is already married!", proposee.DisplayOrUsername())
			return
		}

		ctx.RespondWarningf(
			"%s is already married to %s (@%s)!", proposee.DisplayOrUsername(), spouse.DisplayName, spouse.Username,
		)
		return
	}

	for k, v := range proposals {
		if k == ctx.Interaction.SenderID() || v == ctx.Interaction.SenderID() {
			ctx.RespondWarning("You are already in a proposal!")
			return
		}
		if k == proposeeID || v == proposeeID {
			ctx.RespondWarningf("%s is already in a proposal!", proposee.DisplayOrUsername())
			return
		}
	}

	embed := discord.Embed{
		Color: marriageColour,
		Author: &discord.EmbedAuthor{
			Name: fmt.Sprintf("Proposal to %s", proposee.DisplayOrUsername()),
			Icon: proposeeAvatar,
		},
		Description: fmt.Sprintf(
			"Do you accept %s hand in marriage?",
			util.Possessive(ctx.Interaction.Sender().DisplayName),
		),
		Footer: &discord.EmbedFooter{
			Text: "💍 Proposed",
		},
		Timestamp: discord.Timestamp(ctx.Interaction.ID.Time()),
	}

	msg := fmt.Sprintf("%s has proposed to you, %s!", ctx.Interaction.Sender().DisplayName, proposeeID.Mention())

	err = ctx.RespondMessage(api.InteractionResponseData{
		Content:    option.NewNullableString(msg),
		Embeds:     &[]discord.Embed{embed},
		Components: discord.ComponentsPtr(proposalButtons),
	})

	if err != nil {
		ctx.RespondErrorf("Error occurred while proposing to %s.", proposee.DisplayOrUsername())
		return
	}

	proposals[ctx.Interaction.SenderID()] = proposeeID

	ch, cancel := ctx.State.ChanFor(func(ev any) bool {
		if i, ok := ev.(*gateway.InteractionCreateEvent); ok {
			if _, ok := i.Data.(*discord.ButtonInteraction); ok {
				if i.Message.Interaction.ID == ctx.Interaction.ID && i.SenderID() == proposeeID {
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
				if i.Message.Interaction.ID == ctx.Interaction.ID && i.SenderID() == proposeeID {
					itx := &router.InteractionCtx{
						Router:      ctx.Router,
						Interaction: &i.InteractionEvent,
					}

					if i.Data.(*discord.ButtonInteraction).CustomID == "YES" {
						proposalAccepted(ctx, itx, proposee)
					}
					proposalRejected(ctx, itx, proposee.DisplayOrUsername())
				}
			}
		}
		disableProposalButtons(ctx)
		cancel()

	case <-time.After(24 * time.Hour):
		disableProposalButtons(ctx)
		delete(proposals, ctx.Interaction.SenderID())
		cancel()
	}
}

func proposalAccepted(ctx router.CommandCtx, itx *router.InteractionCtx, proposee discord.User) {
	proposalFound := false
	for k, v := range proposals {
		if k == ctx.Interaction.SenderID() && v == proposee.ID ||
			k == proposee.ID && v == ctx.Interaction.SenderID() {

			proposalFound = true
			break
		}
	}

	if !proposalFound {
		itx.RespondWarning("This proposal is no longer valid!")
		return
	}

	added, err := db.Marriages.Add(ctx.Interaction.SenderID(), proposee.ID)
	if err != nil {
		log.Println(err)
		itx.RespondError("Error occurred while accepting the proposal.")
		return
	}
	if !added {
		itx.RespondWarningf("You and %s are already married!", proposee.DisplayName)
		return
	}

	err = itx.RespondSimple(
		fmt.Sprintf(
			"Congratulations! You and %s are now married! 💗",
			ctx.Interaction.Sender().Mention(),
		),
		discord.Embed{
			Author: &discord.EmbedAuthor{
				Name: fmt.Sprintf(
					"%s and %s Wedding",
					ctx.Interaction.Sender().DisplayName,
					util.Possessive(proposee.DisplayName),
				),
			},
			Description: fmt.Sprintf(
				"💒 %s and %s were married %s 💐",
				ctx.Interaction.Sender().DisplayName,
				proposee.DisplayName,
				dctools.Timestamp(ctx.Interaction.ID.Time()),
			),
			Footer: &discord.EmbedFooter{
				Text: "💍 Proposed",
			},
			Timestamp: discord.Timestamp(ctx.Interaction.ID.Time()),
			Color:     marriageColour,
		},
	)

	if err != nil {
		log.Println(err)
		itx.RespondError("Error occurred while accepting the proposal.")
		return
	}

	delete(proposals, ctx.Interaction.SenderID())
}

func proposalRejected(ctx router.CommandCtx, itx *router.InteractionCtx, proposeeName string) {
	itx.RespondTextf("%s rejected the marriage proposal.", proposeeName)
	delete(proposals, ctx.Interaction.SenderID())
}

func disableProposalButtons(ctx router.CommandCtx) {
	d := dctools.DisabledButtons(*proposalButtons)
	ctx.State.EditInteractionResponse(
		ctx.Interaction.AppID,
		ctx.Interaction.Token,
		api.EditInteractionResponseData{
			Components: discord.ComponentsPtr(&d),
		},
	)
}
