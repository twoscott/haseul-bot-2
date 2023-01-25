package notifications

import (
	"fmt"
	"log"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var notificationsDeleteCommand = &router.SubCommand{
	Name:        "delete",
	Description: "Removes a keyword notification",
	Handler: &router.CommandHandler{
		Executor:      notificationsDeleteExec,
		Autocompleter: notificationKeywordCompleter,
		Ephemeral:     true,
	},
	Options: []discord.CommandOptionValue{
		&discord.StringOption{
			OptionName:   "keyword",
			Description:  "The keyword to remove",
			Required:     true,
			Autocomplete: true,
		},
		&discord.IntegerOption{
			OptionName:  "scope",
			Description: "Where to delete the keyword from",
			Choices: []discord.IntegerChoice{
				{Name: "Server", Value: serverScope},
				{Name: "Global", Value: globalScope},
			},
		},
	},
}

func notificationsDeleteExec(ctx router.CommandCtx) {
	rawKeyword := ctx.Options.Find("keyword").String()
	keyword := strings.ToLower(rawKeyword)

	scope, _ := ctx.Options.Find("scope").IntValue()

	switch scope {
	case serverScope:
		removeServerNoti(ctx, keyword)
	case globalScope:
		removeGlobalNoti(ctx, keyword)
	default:
		ctx.RespondError("Invalid notification scope selected.")
	}
}

func removeServerNoti(ctx router.CommandCtx, keyword string) {
	ok, err := db.Notifications.Remove(
		keyword, ctx.Interaction.SenderID(), ctx.Interaction.GuildID,
	)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while removing keyword from the database.",
		)
		return
	}
	if !ok {
		ctx.RespondWarning(
			"This keyword is not in your server notifications list.",
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

	ctx.RespondSuccess(
		"Notification was removed successfully.",
	)

	var guildName string
	guild, err := ctx.State.Guild(ctx.Interaction.GuildID)
	if err != nil {
		guildName = "the server"
	} else {
		guildName = guild.Name
	}

	dmMsg := fmt.Sprintf(
		"You will no longer be notified when '%s' is mentioned in %s",
		keyword, guildName,
	)

	ctx.State.SendMessage(dmChannel.ID, dmMsg)
}

func removeGlobalNoti(ctx router.CommandCtx, keyword string) {
	ok, err := db.Notifications.RemoveGlobal(
		keyword, ctx.Interaction.SenderID(),
	)
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while removing keyword from the database.",
		)
		return
	}
	if !ok {
		ctx.RespondWarning(
			"This keyword is not in your global notifications list.",
		)
		return
	}

	ctx.RespondSuccess(
		"Notification was removed successfully.",
	)

	dmChannel, err := ctx.State.CreatePrivateChannel(ctx.Interaction.SenderID())
	if err != nil {
		log.Println(err)
		ctx.RespondError(
			"Error occurred while trying to DM you.",
		)
		return
	}

	dmMsg := fmt.Sprintf(
		"You will no longer be notified when '%s' is mentioned globally",
		keyword,
	)

	ctx.State.SendMessage(dmChannel.ID, dmMsg)
}

func notificationKeywordCompleter(ctx router.AutocompleteCtx) {
	keyword := ctx.Options.Find("keyword").String()

	notis, err := db.Notifications.GetByGuildAndGlobalUser(
		ctx.Interaction.SenderID(), ctx.Interaction.GuildID,
	)
	if err != nil {
		log.Println(err)
		ctx.RespondChoices(nil)
		return
	}
	if len(notis) == 0 {
		ctx.RespondChoices(nil)
		return
	}

	keywords := make([]string, 0, len(notis))
	for _, n := range notis {
		keywords = append(keywords, n.Keyword)
	}

	var choices api.AutocompleteStringChoices
	if keyword == "" {
		choices = dctools.MakeStringChoices(keywords)
	} else {
		matches := util.SearchSort(keywords, keyword)
		choices = dctools.MakeStringChoices(matches)
	}

	ctx.RespondChoices(choices)
}
