package logs

import (
	"log"
	"strconv"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
)

var logsWelcomeTitleCommand = &router.SubCommand{
	Name:        "title",
	Description: "Edit the new member welcome title",
	Handler: &router.CommandHandler{
		Executor:     logsWelcomeTitleExec,
		ModalHandler: logsWelcomeEditModalSubmit,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:  "title",
			Description: "The title of the welcome message for new members",
			Required:    true,
			MaxLength:   option.NewInt(32),
		},
	},
}

func logsWelcomeTitleExec(ctx router.CommandCtx) {
	// title := ctx.Options.Find("title").String()

	// _, err := db.Guilds.SetWelcomeTitle(ctx.Interaction.GuildID, title)
	// if err != nil {
	// 	log.Println(err)
	// 	ctx.RespondError("Error occurred while setting welcome title.")
	// 	return
	// }

	// // ctx.RespondSuccess("Welcome title edited.")
	// err = ctx.Respond(api.InteractionResponseData{
	// 	Content: option.NewNullableString("Welcome title edited"),
	// 	Components: discord.ComponentsPtr(
	// 		&discord.ActionRowComponent{
	// 			&discord.StringSelectComponent{
	// 				CustomID:    "test1",
	// 				ValueLimits: [2]int{1, 4},
	// 				Options: []discord.SelectOption{
	// 					{Label: "Test 1", Value: "1"},
	// 					{Label: "Test 2", Value: "2"},
	// 					{Label: "Test 3", Value: "3"},
	// 					{Label: "Test 4", Value: "4"},
	// 				},
	// 			},
	// 		},
	// 		&discord.ActionRowComponent{
	// 			&discord.ChannelSelectComponent{
	// 				CustomID:    "test2",
	// 				ValueLimits: [2]int{1, 25},
	// 				ChannelTypes: []discord.ChannelType{
	// 					discord.GuildText,
	// 				},
	// 			},
	// 		},
	// 		&discord.ActionRowComponent{
	// 			&discord.UserSelectComponent{
	// 				CustomID:    "test3",
	// 				ValueLimits: [2]int{1, 25},
	// 			},
	// 		},
	// 	),
	// })

	// log.Println(err)

	customID := strconv.FormatInt(int64(ctx.Interaction.ID), 10)

	welcomeConfig, err := db.Guilds.WelcomeConfig(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred while fetching welcome config.")
		return
	}

	log.Println("editing welcome")

	defaultTitle := welcomeConfig.Title()
	defaultMessage := welcomeConfig.Message()
	defaultColour := welcomeConfig.Colour().String()
	ctx.RespondWithModal(api.InteractionResponseData{
		CustomID: option.NewNullableString(customID),
		Title:    option.NewNullableString("Edit Welcome Message"),
		Components: discord.ComponentsPtr(
			&discord.TextInputComponent{
				CustomID:     "title",
				Style:        discord.TextInputShortStyle,
				Label:        "Title",
				Value:        defaultTitle,
				LengthLimits: [2]int{0, 32},
			},
			&discord.TextInputComponent{
				CustomID:     "message",
				Style:        discord.TextInputParagraphStyle,
				Label:        "Message",
				Value:        defaultMessage,
				LengthLimits: [2]int{0, 1024},
			},
			&discord.TextInputComponent{
				CustomID:     "colour",
				Style:        discord.TextInputShortStyle,
				Label:        "Hex Colour",
				Value:        defaultColour,
				LengthLimits: [2]int{6, 7},
			},
		),
	})
}

func logsWelcomeEditModalSubmit(ctx router.ModalCtx) {
	log.Println("Modal response received")

	titleCmp := ctx.Components.Find("title").(*discord.TextInputComponent)
	messageCmp := ctx.Components.Find("message").(*discord.TextInputComponent)
	colourCmp := ctx.Components.Find("colour").(*discord.TextInputComponent)

	log.Println(titleCmp.Value)
	log.Println(messageCmp.Value)
	log.Println(colourCmp.Value)
}
