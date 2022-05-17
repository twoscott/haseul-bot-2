package misc

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/twoscott/haseul-bot-2/router"
)

var testCommand = &router.Command{
	Name:      "test",
	UseTyping: true,
	Run:       testRun,
}

func testRun(ctx router.CommandCtx, args []string) {
	msg, _ := ctx.State.SendMessageComplex(ctx.Msg.ChannelID, api.SendMessageData{
		Content: "Test",
		Components: []discord.Component{
			discord.ActionRowComponent{
				Components: []discord.Component{
					discord.ButtonComponent{
						Label:    "Hello World!",
						CustomID: "first_button",
						Emoji: &discord.ButtonEmoji{
							Name: "👋",
						},
						Style: discord.PrimaryButton,
					},
					discord.ButtonComponent{
						Label:    "Secondary",
						CustomID: "second_button",
						Style:    discord.SecondaryButton,
					},
					discord.ButtonComponent{
						Label:    "Success",
						CustomID: "success_button",
						Style:    discord.SuccessButton,
					},
					discord.ButtonComponent{
						Label:    "Danger",
						CustomID: "danger_button",
						Style:    discord.DangerButton,
					},
					discord.ButtonComponent{
						Label: "Link",
						URL:   "https://google.com",
						Style: discord.LinkButton,
					},
				},
			},
		},
	})

	<-time.After(time.Second)

	c := []discord.Component{
		discord.ActionRowComponent{
			Components: []discord.Component{
				discord.ButtonComponent{
					Label:    "Hello World!",
					CustomID: "first_button",
					Emoji: &discord.ButtonEmoji{
						Name: "👋",
					},
					Style: discord.PrimaryButton,
				},
				discord.ButtonComponent{
					Label:    "Secondary",
					CustomID: "second_button",
					Style:    discord.SecondaryButton,
					Disabled: true,
				},
				discord.ButtonComponent{
					Label:    "Success",
					CustomID: "success_button",
					Style:    discord.SuccessButton,
				},
				discord.ButtonComponent{
					Label:    "Danger",
					CustomID: "danger_button",
					Style:    discord.DangerButton,
				},
				discord.ButtonComponent{
					Label: "Link",
					URL:   "https://google.com",
					Style: discord.LinkButton,
				},
			},
		},
	}

	t := api.EditMessageData{
		Content:    option.NewNullableString("A"),
		Components: &c,
	}

	var indent bytes.Buffer
	b, _ := json.Marshal(t)
	json.Indent(&indent, b, "", "    ")
	log.Println(indent.String())

	_, err := ctx.State.EditMessageComplex(msg.ChannelID, msg.ID, t)
	log.Println(err)
}
