package cmdutil

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

func ReplyWithPaging(
	ctx router.CommandCtx,
	msg *discord.Message,
	messagePages []router.MessagePage) (*discord.Message, error) {

	return replyWithPaging(ctx, msg, messagePages, false)
}

func ReplyWithConfirmationPaging(
	ctx router.CommandCtx,
	msg *discord.Message,
	messagePages []router.MessagePage) (*discord.Message, error) {

	return replyWithPaging(ctx, msg, messagePages, true)
}

func replyWithPaging(
	ctx router.CommandCtx,
	msg *discord.Message,
	messagePages []router.MessagePage,
	confirmation bool) (*discord.Message, error) {

	if len(messagePages) < 2 {
		return dctools.ReplyNoPing(ctx.State, msg,
			messagePages[0].Content,
			messagePages[0].Embeds...,
		)
	}

	var (
		sentMsg *discord.Message
		err     error
	)
	if confirmation {
		sentMsg, err = dctools.ReplyWithConfirmationPager(ctx.State, msg,
			messagePages[0].Content,
			messagePages[0].Embeds...,
		)
	} else {
		sentMsg, err = dctools.ReplyWithMessagePager(ctx.State, msg,
			messagePages[0].Content,
			messagePages[0].Embeds...,
		)
	}
	if err != nil {
		return sentMsg, err
	}

	err = ctx.AddButtonPager(router.ButtonPagerOptions{
		AuthorID:  ctx.Msg.Author.ID,
		ChannelID: sentMsg.ChannelID,
		MessageID: sentMsg.ID,
		Pages:     messagePages,
	})

	return sentMsg, err
}

func SendWithPaging(
	ctx router.CommandCtx,
	channelID discord.ChannelID,
	messagePages []router.MessagePage) (*discord.Message, error) {

	return sendWithPaging(ctx, channelID, messagePages, false)
}

func SendWithConfirmationPaging(
	ctx router.CommandCtx,
	channelID discord.ChannelID,
	messagePages []router.MessagePage) (*discord.Message, error) {

	return sendWithPaging(ctx, channelID, messagePages, true)
}

func sendWithPaging(
	ctx router.CommandCtx,
	channelID discord.ChannelID,
	messagePages []router.MessagePage,
	confirmation bool) (*discord.Message, error) {

	if len(messagePages) < 2 {
		return ctx.State.SendMessage(channelID,
			messagePages[0].Content,
			messagePages[0].Embeds...,
		)
	}

	var (
		sentMsg *discord.Message
		err     error
	)
	if confirmation {
		sentMsg, err = dctools.SendWithConfirmationPager(ctx.State, channelID,
			messagePages[0].Content,
			messagePages[0].Embeds...,
		)
	} else {
		sentMsg, err = dctools.SendWithMessagePager(ctx.State, channelID,
			messagePages[0].Content,
			messagePages[0].Embeds...,
		)
	}
	if err != nil {
		return sentMsg, err
	}

	err = ctx.AddButtonPager(router.ButtonPagerOptions{
		AuthorID:  ctx.Msg.Author.ID,
		ChannelID: sentMsg.ChannelID,
		MessageID: sentMsg.ID,
		Pages:     messagePages,
	})

	return sentMsg, err
}
