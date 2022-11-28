package modules

import (
	"github.com/twoscott/haseul-bot-2/modules/bot"
	"github.com/twoscott/haseul-bot-2/modules/command"
	"github.com/twoscott/haseul-bot-2/modules/emoji"
	"github.com/twoscott/haseul-bot-2/modules/lastfm"
	"github.com/twoscott/haseul-bot-2/modules/logs"
	"github.com/twoscott/haseul-bot-2/modules/misc"
	"github.com/twoscott/haseul-bot-2/modules/notification"
	"github.com/twoscott/haseul-bot-2/modules/server"
	"github.com/twoscott/haseul-bot-2/modules/twitter"
	"github.com/twoscott/haseul-bot-2/modules/user"
	"github.com/twoscott/haseul-bot-2/modules/vlive"
	"github.com/twoscott/haseul-bot-2/modules/youtube"
	"github.com/twoscott/haseul-bot-2/router"
)

func Init(rt *router.Router) {
	bot.Init(rt)
	command.Init(rt)
	emoji.Init(rt)
	lastfm.Init(rt)
	logs.Init(rt)
	misc.Init(rt)
	notification.Init(rt)
	server.Init(rt)
	twitter.Init(rt)
	user.Init(rt)
	vlive.Init(rt)
	youtube.Init(rt)
}
