package message

import (
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

var editMessageCommand = &router.Command{
	Name: "Edit",
	Type: discord.MessageCommand,
	Handler: &router.CommandHandler{
		Executor:  editMessageExec,
		Ephemeral: true,
	},
	RequiredPermissions: discord.NewPermissions(
		discord.PermissionManageMessages,
	),
}

func editMessageExec(ctx router.CommandCtx) {
	var permissions discord.Permissions
	bot, err := ctx.State.Me()
	if err == nil {
		permissions, err = ctx.State.Permissions(
			ctx.Interaction.ChannelID,
			bot.ID,
		)
	}
	if err == nil && !dctools.HasAnyPermOrAdmin(
		permissions, discord.PermissionViewChannel) {

		ctx.RespondWarning(
			"I do not have permission to view this channel.",
		)
		return
	}

	msg, err := ctx.State.Message(
		ctx.Interaction.ChannelID,
		ctx.Command.TargetMessageID(),
	)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching message data.")
		return
	}

	if msg.Author.ID != bot.ID {
		ctx.RespondWarning("This message was not sent by me.")
		return
	}

	if msg.Interaction != nil {
		ctx.RespondWarning("This message was sent in response to a command.")
		return
	}

	textBox := &discord.TextInputComponent{
		CustomID:     "CONTENT",
		Label:        "Content",
		Style:        discord.TextInputParagraphStyle,
		Placeholder:  "Enter the new message content to edit the message with.",
		Value:        msg.Content,
		LengthLimits: [2]int{0, 2000},
	}

	customID := ctx.Interaction.ID.String()

	evCh, cancel := ctx.State.ChanFor(func(ev interface{}) bool {
		i, ok := ev.(*gateway.InteractionCreateEvent)
		if !ok {
			return ok
		}

		m, ok := i.Data.(*discord.ModalInteraction)
		return ok && m.CustomID == discord.ComponentID(customID)
	})
	defer cancel()

	err = ctx.RespondModal(
		api.InteractionResponseData{
			CustomID:   option.NewNullableString(customID),
			Title:      option.NewNullableString("Edit Message"),
			Components: discord.ComponentsPtr(textBox),
		},
	)
	if err != nil {
		log.Println(err)
		ctx.RespondGenericError()
		return
	}

	var (
		submit *discord.InteractionEvent
		modal  *discord.ModalInteraction
	)
	select {
	case ev := <-evCh:
		icv := ev.(*gateway.InteractionCreateEvent)
		submit = &icv.InteractionEvent
		modal = submit.Data.(*discord.ModalInteraction)
	case <-time.After(30 * time.Minute):
		ctx.RespondWarning("Modal timed out.")
		return
	}
	if submit == nil || modal == nil {
		ctx.RespondGenericError()
		return
	}

	modalCtx := router.ModalCtx{
		InteractionCtx: &router.InteractionCtx{
			Router:      ctx.Router,
			Interaction: submit,
			Ephemeral:   true,
		},
		Modal: modal,
	}

	processModalSubmit(modalCtx, msg)
}

func processModalSubmit(ctx router.ModalCtx, msg *discord.Message) {
	component := ctx.Modal.Components.Find("CONTENT")
	if component == nil {
		ctx.RespondError("Error occurred fetching edited message content.")
		return
	}

	contentComponent, ok := component.(*discord.TextInputComponent)
	if !ok {
		ctx.RespondError("Error occurred processing submitted modal data.")
		return
	}

	if contentComponent.Value == "" && len(msg.Attachments) == 0 {
		ctx.RespondWarning(
			"This message would be deleted if the content were removed.",
		)
		return
	}

	_, err := ctx.State.EditMessageComplex(ctx.Interaction.ChannelID, msg.ID,
		api.EditMessageData{
			Content: option.NewNullableString(contentComponent.Value),
		},
	)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while editing message.")
		return
	}

	ctx.RespondSuccess(
		fmt.Sprintf("Message edited successfully. %s", msg.URL()),
	)
}
