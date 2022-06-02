package notifications

import (
	"fmt"
	"log"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/database/notifdb"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

const notificationLimit = 25

var notiAddCommand = &router.SubCommand{
	Name:        "add",
	Description: "Adds a keyword notification",
	Handler: &router.CommandHandler{
		Executor:  notiAddExec,
		Ephemeral: true,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:  "keyword",
			Description: "The keyword to be notified for mentions of",
			Required:    true,
		},
		&discord.IntegerOption{
			OptionName:  "scope",
			Description: "Where to be notified for the keyword",
			Choices: []discord.IntegerChoice{
				{Name: "Server", Value: int(serverScope)},
				{Name: "Global", Value: int(globalScope)},
			},
		},
		&discord.IntegerOption{
			OptionName:  "type",
			Description: "How to be notified of the keyword",
			Choices: []discord.IntegerChoice{
				{Name: "Normal", Value: int(notifdb.NormalNotification)},
				{Name: "Strict", Value: int(notifdb.StrictNotification)},
				{Name: "Lenient", Value: int(notifdb.LenientNotification)},
				{Name: "Anarchy", Value: int(notifdb.AnarchyNotification)},
			},
		},
	},
}

func notiAddExec(ctx router.CommandCtx) {
	rawKeyword := ctx.Options.Find("keyword").String()
	keyword := strings.ToLower(rawKeyword)
	if keyword == "" {
		ctx.RespondWarning(
			"Please provide a keyword to get notified for.",
		)
		return
	}
	if len([]rune(keyword)) > 128 {
		ctx.RespondWarning(
			"Keywords must be less than 128 characters in length.",
		)
		return
	}

	keywordScope, _ := ctx.Options.Find("scope").IntValue()
	typeOption, _ := ctx.Options.Find("type").IntValue()
	keywordType := notifdb.NotificationType(typeOption)

	switch keywordScope {
	case serverScope:
		addServerNoti(ctx, keyword, keywordType)
	case globalScope:
		addGlobalNoti(ctx, keyword, keywordType)
	default:
		ctx.RespondError("Invalid notification scope selected.")
	}
}

func addServerNoti(
	ctx router.CommandCtx,
	keyword string,
	keywordType notifdb.NotificationType) {

	notifications, err := db.Notifications.GetByGuildUser(
		ctx.Interaction.SenderID(), ctx.Interaction.GuildID,
	)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while checking your notifications.",
		)
		return
	}
	if len(notifications) >= 25 {
		ctx.RespondWarning(
			"You cannot have more than 25 notifications set up in a server. " +
				"You may remove server notifications and re-add them " +
				"as global notifications.",
		)
		return
	}

	ok, err := db.Notifications.Add(
		keyword,
		ctx.Interaction.SenderID(),
		keywordType,
		ctx.Interaction.GuildID,
	)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while adding keyword to the database.",
		)
		return
	}
	if !ok {
		ctx.RespondWarning(
			"You are already notified of this keyword.",
		)
		return
	}

	dmChannel, err := ctx.State.CreatePrivateChannel(ctx.Interaction.SenderID())
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while trying to DM you.",
		)
		return
	}

	var guildName string
	guild, err := ctx.State.Guild(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		guildName = "the server"
	} else {
		guildName = guild.Name
	}

	dmMsg := fmt.Sprintf(
		"You will now be notified when '%s' is mentioned in %s.",
		keyword, guildName,
	)

	_, err = ctx.State.SendMessage(dmChannel.ID, dmMsg)
	if dctools.ErrCannotDM(err) {
		ctx.RespondWarning(
			"I am unable to DM you. " +
				"Please open your DMs to server members in your settings.",
		)
		db.Notifications.Remove(
			keyword, ctx.Interaction.SenderID(), ctx.Interaction.GuildID,
		)
		return
	}

	ctx.RespondSuccess(
		"Notification was added successfully.",
	)
}

func addGlobalNoti(
	ctx router.CommandCtx,
	keyword string,
	keywordType notifdb.NotificationType) {

	notifications, err := db.Notifications.GetByUser(ctx.Interaction.SenderID())
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while checking your notifications.",
		)
		return
	}
	if len(notifications) >= 25 {
		ctx.RespondWarning(
			"You cannot have more than 25 global notifications set up.",
		)
		return
	}

	ok, err := db.Notifications.AddGlobal(
		keyword,
		ctx.Interaction.SenderID(),
		keywordType,
	)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while adding keyword to the database.",
		)
		return
	}
	if !ok {
		ctx.RespondWarning(
			"You are already notified of this keyword.",
		)
		return
	}

	dmChannel, err := ctx.State.CreatePrivateChannel(ctx.Interaction.SenderID())
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while trying to DM you.",
		)
		return
	}

	dmMsg := fmt.Sprintf(
		"You will now be notified when '%s' is mentioned in globally.",
		keyword,
	)

	_, err = ctx.State.SendMessage(dmChannel.ID, dmMsg)
	if dctools.ErrCannotDM(err) {
		ctx.RespondWarning(
			"I am unable to DM you. " +
				"Please open your DMs to server members in your settings.",
		)
		db.Notifications.RemoveGlobal(
			keyword, ctx.Interaction.SenderID(),
		)
		return
	}

	ctx.RespondSuccess(
		"Notification was added successfully.",
	)
}
