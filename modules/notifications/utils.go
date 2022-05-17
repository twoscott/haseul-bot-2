package notifications

import (
	"regexp"
	"strings"

	"github.com/twoscott/haseul-bot-2/database/notifdb"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

var messageCleanRegex = regexp.MustCompile("([\\*\\`\\~\\_\\]\\)])")

func cleanMessageContent(content string) string {
	return messageCleanRegex.ReplaceAllString(content, "")
}

func getKeyword(
	ctx router.CommandCtx, args []string) (string, notifdb.NotificationType) {

	var keyType notifdb.NotificationType
	keywordStart := 1

	switch strings.ToLower(args[0]) {
	default:
		keywordStart = 0
		fallthrough
	case "normal":
		keyType = notifdb.NormalNotification
	case "strict":
		keyType = notifdb.StrictNotification
	case "lenient":
		keyType = notifdb.LenientNotification
	case "anarchy":
		keyType = notifdb.AnarchyNotification
	}

	rawKeyword := util.TrimArgs(ctx.Msg.Content, ctx.Length+keywordStart)
	keyword := strings.ToLower(rawKeyword)

	return keyword, keyType
}
