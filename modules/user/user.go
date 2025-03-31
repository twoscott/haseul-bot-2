package user

import "github.com/twoscott/haseul-bot-2/router"

const (
	serverType int64 = iota
	globalType
)

var userCommand = &router.Command{
	Name:        "user",
	Description: "Commands pertaining to Discord users",
}
