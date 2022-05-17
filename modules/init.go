package modules

import (
	"github.com/twoscott/haseul-bot-2/modules/emoji"
	"github.com/twoscott/haseul-bot-2/modules/information"
	"github.com/twoscott/haseul-bot-2/modules/lastfm"
	"github.com/twoscott/haseul-bot-2/modules/misc"
	"github.com/twoscott/haseul-bot-2/modules/notifications"
	"github.com/twoscott/haseul-bot-2/modules/twitter"
	"github.com/twoscott/haseul-bot-2/modules/vlive"
	"github.com/twoscott/haseul-bot-2/modules/youtube"
	"github.com/twoscott/haseul-bot-2/router"
)

func Init(rt *router.Router) {
	emoji.Init(rt)
	information.Init(rt)
	lastfm.Init(rt)
	misc.Init(rt)
	notifications.Init(rt)
	twitter.Init(rt)
	vlive.Init(rt)
	youtube.Init(rt)
}
